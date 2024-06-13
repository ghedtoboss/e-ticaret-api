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
)

func main() {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		log.Fatal("JWT_SECRET_KEY env is not set")
	}

	db := db.InitDB("root:Eses147852@tcp(127.0.0.1:3306)/e_commerce_api?parseTime=true")
	defer db.Close()
	fmt.Println("Veritabanına bağlanıldı.")

	r := mux.NewRouter()

	appHandler := &handlers.AppHandler{DB: db}

	//routes
	r.Handle("/register", appHandler.Register()).Methods("POST")
	r.Handle("/login", appHandler.Login()).Methods("POST")
	r.Handle("/product", middleware.JWTMiddleware(middleware.RoleMiddleware("seller")(appHandler.AddProduct()))).Methods("POST")
	r.Handle("/product/{id}", middleware.JWTMiddleware(middleware.RoleMiddleware("seller")(appHandler.UpdateProduct()))).Methods("PUT")
	r.Handle("/product/{id}", middleware.JWTMiddleware(middleware.RoleMiddleware("seller")(appHandler.DeleteProduct()))).Methods("DELETE")
	r.Handle("/products", appHandler.GetProducts()).Methods("GET")
	r.Handle("/cart", middleware.JWTMiddleware(appHandler.AddToCart())).Methods("POST")
	r.Handle("/cart", middleware.JWTMiddleware(appHandler.GetCartItems())).Methods("GET")
	r.Handle("/carts/remove/{item_id}", middleware.JWTMiddleware(appHandler.RemoveFromCart())).Methods("DELETE")
	r.Handle("/carts/decrease/{item_id}", middleware.JWTMiddleware(appHandler.DecreaseItemQuantity())).Methods("PUT")
	r.Handle("/carts/increase/{item_id}", middleware.JWTMiddleware(appHandler.IncreaseItemQuantity())).Methods("PUT")
	r.Handle("/carts/remove/cart/items", middleware.JWTMiddleware(appHandler.RemoveCartItems())).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
