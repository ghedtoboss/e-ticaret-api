package handlers

import (
	"e-ticaret-api/models"
	"encoding/json"
	"net/http"
)

// GetUsers godoc
// @Summary Get all users
// @Description Get all users by admin
// @Tags admin
// @Accept  json
// @Produce  json
// @Success 200 {array} models.User
// @Failure 403 {string} string "Only admin can access this endpoint"
// @Failure 500 {string} string "Error fetching users"
// @Router /admin/users [get]
// @Security ApiKeyAuth
func (db *AppHandler) GetUsers() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userRole := r.Context().Value("role").(string)
		if userRole != "admin" {
			http.Error(w, "Only admin can access this endpoint", http.StatusForbidden)
			return
		}

		rows, err := db.DB.Query("SELECT id, email, name, role FROM users")
		if err != nil {
			http.Error(w, "Error fetching users", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []models.User
		for rows.Next() {
			var user models.User
			if err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.Role); err != nil {
				http.Error(w, "Error scanning user", http.StatusInternalServerError)
				return
			}
			users = append(users, user)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	})
}

// AdminAddProduct godoc
// @Summary Add a product by admin
// @Description Add a new product by admin
// @Tags admin
// @Accept  json
// @Produce  json
// @Param   product  body     models.Product  true  "Product"
// @Success 201 {object} models.Product
// @Failure 403 {string} string "Only admin can add products"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Error adding product"
// @Router /admin/products [post]
// @Security ApiKeyAuth
func (db *AppHandler) AdminAddProduct() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userRole := r.Context().Value("role").(string)
		if userRole != "admin" {
			http.Error(w, "Only admin can add products", http.StatusForbidden)
			return
		}

		var product models.Product
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		_, err := db.DB.Exec("INSERT INTO products (name, description, quantity, price, seller_id, category, image_url) VALUES (?, ?, ?, ?, ?, ?, ?)",
			product.Name, product.Description, product.Quantity, product.Price, product.SellerID, product.Category, product.ImageURL)
		if err != nil {
			http.Error(w, "Error adding product", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(product)
	})
}

// GetAllOrders godoc
// @Summary Get all orders by admin
// @Description Get all orders by admin
// @Tags admin
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Order
// @Failure 403 {string} string "Only admin can access this endpoint"
// @Failure 500 {string} string "Error fetching orders"
// @Router /admin/orders [get]
// @Security ApiKeyAuth
func (db *AppHandler) GetAllOrders() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userRole := r.Context().Value("role").(string)
		if userRole != "admin" {
			http.Error(w, "Only admin can access this endpoint", http.StatusForbidden)
			return
		}

		rows, err := db.DB.Query("SELECT id, user_id, total_price, created_at, status FROM orders")
		if err != nil {
			http.Error(w, "Error fetching orders", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var orders []models.Order
		for rows.Next() {
			var order models.Order
			if err := rows.Scan(&order.ID, &order.UserID, &order.TotalPrice, &order.CreatedAt, &order.Status); err != nil {
				http.Error(w, "Error scanning order", http.StatusInternalServerError)
				return
			}
			orders = append(orders, order)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orders)
	})
}
