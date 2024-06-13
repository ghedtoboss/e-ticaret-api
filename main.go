package main

import (
	"e-ticaret-api/db"
	"e-ticaret-api/handlers"
	"e-ticaret-api/middleware"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger" // swagger

	_ "e-ticaret-api/docs" // Swag dokümantasyonu için gerekli
)

// @title E-Ticaret API
// @version 1.0
// @description E-Ticaret API dokümantasyonu.

// @host localhost:8080
// @BasePath /
// @schemes http

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		log.Fatal("JWT_SECRET_KEY env is not set")
	}

	databaseUrl := os.Getenv("DATABASE_URL")
	db := db.InitDB(databaseUrl)
	defer db.Close()
	fmt.Println("Veritabanına bağlanıldı.")

	r := mux.NewRouter()

	appHandler := &handlers.AppHandler{DB: db}

	// Swagger endpoint
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	//routes
	// @Summary Register a new user
	// @Description Register a new user with email, password, name, and role
	// @Tags auth
	// @Accept  json
	// @Produce  json
	// @Param   user     body     models.User     true  "User"
	// @Success 201 {object} models.User
	// @Failure 400 {string} string "Invalid request"
	// @Failure 500 {string} string "Internal server error"
	// @Router /register [post]
	r.Handle("/register", appHandler.Register()).Methods("POST")

	// @Summary Login a user
	// @Description Login a user with email and password
	// @Tags auth
	// @Accept  json
	// @Produce  json
	// @Param   user     body     models.User     true  "User"
	// @Success 200 {object} map[string]string
	// @Failure 400 {string} string "Invalid request"
	// @Failure 401 {string} string "Unauthorized"
	// @Router /login [post]
	r.Handle("/login", appHandler.Login()).Methods("POST")

	// @Summary Add a new product
	// @Description Add a new product by seller
	// @Tags products
	// @Accept  json
	// @Produce  json
	// @Param   product  body     models.Product  true  "Product"
	// @Success 201 {object} models.Product
	// @Failure 400 {string} string "Invalid request"
	// @Failure 500 {string} string "Internal server error"
	// @Router /product [post]
	// @Security ApiKeyAuth
	r.Handle("/product", middleware.JWTMiddleware(middleware.RoleMiddleware("seller")(appHandler.AddProduct()))).Methods("POST")

	// @Summary Update a product
	// @Description Update a product by seller
	// @Tags products
	// @Accept  json
	// @Produce  json
	// @Param   id       path     int             true  "Product ID"
	// @Param   product  body     models.Product  true  "Product"
	// @Success 200 {object} models.Product
	// @Failure 400 {string} string "Invalid request"
	// @Failure 404 {string} string "Product not found"
	// @Failure 500 {string} string "Internal server error"
	// @Router /product/{id} [put]
	// @Security ApiKeyAuth
	r.Handle("/product/{id}", middleware.JWTMiddleware(middleware.RoleMiddleware("seller")(appHandler.UpdateProduct()))).Methods("PUT")

	// @Summary Delete a product
	// @Description Delete a product by seller
	// @Tags products
	// @Accept  json
	// @Produce  json
	// @Param   id  path  int  true  "Product ID"
	// @Success 200 {string} string "Product deleted"
	// @Failure 400 {string} string "Invalid request"
	// @Failure 404 {string} string "Product not found"
	// @Failure 500 {string} string "Internal server error"
	// @Router /product/{id} [delete]
	// @Security ApiKeyAuth
	r.Handle("/product/{id}", middleware.JWTMiddleware(middleware.RoleMiddleware("seller")(appHandler.DeleteProduct()))).Methods("DELETE")

	// @Summary Get all products
	// @Description Get all products with optional filters
	// @Tags products
	// @Accept  json
	// @Produce  json
	// @Param   category  query    string  false  "Category"
	// @Param   search    query    string  false  "Search"
	// @Param   sort_by   query    string  false  "Sort by"
	// @Param   order     query    string  false  "Order"
	// @Success 200 {array} models.Product
	// @Failure 500 {string} string "Internal server error"
	// @Router /products [get]
	r.Handle("/products", appHandler.GetProducts()).Methods("GET")

	// @Summary Add to cart
	// @Description Add a product to the cart
	// @Tags cart
	// @Accept  json
	// @Produce  json
	// @Param   cartItem  body     models.CartItem  true  "Cart Item"
	// @Success 201 {object} models.CartItem
	// @Failure 400 {string} string "Invalid request"
	// @Failure 500 {string} string "Internal server error"
	// @Router /cart [post]
	// @Security ApiKeyAuth
	r.Handle("/cart", middleware.JWTMiddleware(appHandler.AddToCart())).Methods("POST")

	// @Summary Get cart items
	// @Description Get all items in the cart
	// @Tags cart
	// @Accept  json
	// @Produce  json
	// @Success 200 {array} models.CartItem
	// @Failure 500 {string} string "Internal server error"
	// @Router /cart [get]
	// @Security ApiKeyAuth
	r.Handle("/cart", middleware.JWTMiddleware(appHandler.GetCartItems())).Methods("GET")

	// @Summary Remove item from cart
	// @Description Remove an item from the cart
	// @Tags cart
	// @Accept  json
	// @Produce  json
	// @Param   item_id  path  int  true  "Item ID"
	// @Success 200 {string} string "Item removed from cart"
	// @Failure 400 {string} string "Invalid request"
	// @Failure 404 {string} string "Item not found"
	// @Failure 500 {string} string "Internal server error"
	// @Router /carts/remove/{item_id} [delete]
	// @Security ApiKeyAuth
	r.Handle("/carts/remove/{item_id}", middleware.JWTMiddleware(appHandler.RemoveFromCart())).Methods("DELETE")

	// @Summary Decrease item quantity
	// @Description Decrease the quantity of an item in the cart
	// @Tags cart
	// @Accept  json
	// @Produce  json
	// @Param   item_id  path  int  true  "Item ID"
	// @Success 200 {string} string "Item quantity decreased"
	// @Failure 400 {string} string "Invalid request"
	// @Failure 404 {string} string "Item not found"
	// @Failure 500 {string} string "Internal server error"
	// @Router /carts/decrease/{item_id} [put]
	// @Security ApiKeyAuth
	r.Handle("/carts/decrease/{item_id}", middleware.JWTMiddleware(appHandler.DecreaseItemQuantity())).Methods("PUT")

	// @Summary Increase item quantity
	// @Description Increase the quantity of an item in the cart
	// @Tags cart
	// @Accept  json
	// @Produce  json
	// @Param   item_id  path  int  true  "Item ID"
	// @Success 200 {string} string "Item quantity increased"
	// @Failure 400 {string} string "Invalid request"
	// @Failure 404 {string} string "Item not found"
	// @Failure 500 {string} string "Internal server error"
	// @Router /carts/increase/{item_id} [put]
	// @Security ApiKeyAuth
	r.Handle("/carts/increase/{item_id}", middleware.JWTMiddleware(appHandler.IncreaseItemQuantity())).Methods("PUT")

	// @Summary Clear cart items
	// @Description Clear all items from the cart
	// @Tags cart
	// @Accept  json
	// @Produce  json
	// @Success 200 {string} string "Cart cleared"
	// @Failure 500 {string} string "Internal server error"
	// @Router /carts/remove/cart/items [delete]
	// @Security ApiKeyAuth
	r.Handle("/carts/remove/cart/items", middleware.JWTMiddleware(appHandler.RemoveCartItems())).Methods("DELETE")

	// @Summary Create order
	// @Description Create a new order from the cart items
	// @Tags orders
	// @Accept  json
	// @Produce  json
	// @Success 201 {object} models.Order
	// @Failure 400 {string} string "Invalid request"
	// @Failure 500 {string} string "Internal server error"
	// @Router /order [post]
	// @Security ApiKeyAuth
	r.Handle("/order", middleware.JWTMiddleware(appHandler.CreateOrder())).Methods("POST")

	// @Summary Get user orders
	// @Description Get all orders for a user
	// @Tags orders
	// @Accept  json
	// @Produce  json
	// @Success 200 {array} models.Order
	// @Failure 500 {string} string "Internal server error"
	// @Router /orders [get]
	// @Security ApiKeyAuth
	r.Handle("/orders", middleware.JWTMiddleware(appHandler.GetOrders())).Methods("GET")

	// @Summary Get order items
	// @Description Get all items for a specific order
	// @Tags orders
	// @Accept  json
	// @Produce  json
	// @Param   order_id  path  int  true  "Order ID"
	// @Success 200 {array} models.OrderItem
	// @Failure 400 {string} string "Invalid request"
	// @Failure 404 {string} string "Order not found"
	// @Failure 500 {string} string "Internal server error"
	// @Router /orders/{order_id} [get]
	// @Security ApiKeyAuth
	r.Handle("/orders/{order_id}", middleware.JWTMiddleware(appHandler.GetOrderItems())).Methods("GET")

	// @Summary Update order status
	// @Description Update the status of an order by admin
	// @Tags orders
	// @Accept  json
	// @Produce  json
	// @Param   order_id  path  int  true  "Order ID"
	// @Param   status    body  string  true  "Order Status"
	// @Success 200 {string} string "Order status updated"
	// @Failure 400 {string} string "Invalid request"
	// @Failure 404 {string} string "Order not found"
	// @Failure 500 {string} string "Internal server error"
	// @Router /orders/{order_id}/status [put]
	// @Security ApiKeyAuth
	r.Handle("/orders/{order_id}/status", middleware.JWTMiddleware(middleware.RoleMiddleware("admin")(appHandler.UpdateOrderStatus()))).Methods("PUT")

	// @Summary Get all users
	// @Description Get all users by admin
	// @Tags admin
	// @Accept  json
	// @Produce  json
	// @Success 200 {array} models.User
	// @Failure 500 {string} string "Internal server error"
	// @Router /admin/users [get]
	// @Security ApiKeyAuth
	r.Handle("/admin/users", middleware.JWTMiddleware(middleware.RoleMiddleware("admin")(appHandler.GetUsers()))).Methods("GET")

	// @Summary Add a product by admin
	// @Description Add a new product by admin
	// @Tags admin
	// @Accept  json
	// @Produce  json
	// @Param   product  body     models.Product  true  "Product"
	// @Success 201 {object} models.Product
	// @Failure 400 {string} string "Invalid request"
	// @Failure 500 {string} string "Internal server error"
	// @Router /admin/products [post]
	// @Security ApiKeyAuth
	r.Handle("/admin/products", middleware.JWTMiddleware(middleware.RoleMiddleware("admin")(appHandler.AdminAddProduct()))).Methods("POST")

	// @Summary Get all orders by admin
	// @Description Get all orders by admin
	// @Tags admin
	// @Accept  json
	// @Produce  json
	// @Success 200 {array} models.Order
	// @Failure 500 {string} string "Internal server error"
	// @Router /admin/orders [get]
	// @Security ApiKeyAuth
	r.Handle("/admin/orders", middleware.JWTMiddleware(middleware.RoleMiddleware("admin")(appHandler.GetAllOrders()))).Methods("GET")

	// @Summary Create a return
	// @Description Create a return for an order
	// @Tags Returns
	// @Accept json
	// @Produce json
	// @Security ApiKeyAuth
	// @Param return body models.Return true "Return"
	// @Success 201 {object} models.Return
	// @Router /returns [post]
	r.Handle("/returns", middleware.JWTMiddleware(appHandler.CreateReturn())).Methods("POST")

	// @Summary Get returns
	// @Description Get all returns for a user
	// @Tags Returns
	// @Accept json
	// @Produce json
	// @Security ApiKeyAuth
	// @Success 200 {array} models.Return
	// @Router /returns [get]
	r.Handle("/returns", middleware.JWTMiddleware(appHandler.GetReturns())).Methods("GET")

	// @Summary Create a review
	// @Description Create a review for a product
	// @Tags Reviews
	// @Accept json
	// @Produce json
	// @Security ApiKeyAuth
	// @Param review body models.Review true "Review"
	// @Success 201 {object} models.Review
	// @Router /reviews [post]
	r.Handle("/reviews", middleware.JWTMiddleware(appHandler.CreateReview())).Methods("POST")

	// @Summary Get reviews for a product
	// @Description Get all reviews for a specific product
	// @Tags Reviews
	// @Accept json
	// @Produce json
	// @Param product_id path int true "Product ID"
	// @Success 200 {array} models.Review
	// @Router /reviews/{product_id} [get]
	r.Handle("/reviews/{product_id}", appHandler.GetReviews()).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
