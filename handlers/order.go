package handlers

import (
	"e-ticaret-api/models"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// CreateOrder godoc
// @Summary Create an order
// @Description Create an order for the authenticated user
// @Tags orders
// @Produce  json
// @Success 201 {object} models.Order
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Cart not found"
// @Failure 500 {string} string "Internal server error"
// @Router /order [post]
func (db *AppHandler) CreateOrder() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)

		var cartID int
		err := db.DB.QueryRow("SELECT id FROM carts WHERE user_id = ?", userID).Scan(&cartID)
		if err != nil {
			http.Error(w, "Cart not found", http.StatusNotFound)
			return
		}

		var totalPrice float64
		err = db.DB.QueryRow("SELECT SUM(price * quantity) FROM cart_items WHERE cart_id = ?", cartID).Scan(&totalPrice)
		if err != nil {
			http.Error(w, "Error calculating total price", http.StatusInternalServerError)
			return
		}

		order := models.Order{
			UserID:     userID,
			TotalPrice: totalPrice,
			CreatedAt:  time.Now(),
		}

		tx, err := db.DB.Begin()
		if err != nil {
			http.Error(w, "Transaction begin error", http.StatusInternalServerError)
			return
		}

		res, err := tx.Exec("INSERT INTO orders (user_id, total_price, created_at) VALUES (?, ?, ?)", order.UserID, order.TotalPrice, order.CreatedAt)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error inserting order", http.StatusInternalServerError)
			return
		}

		lastInsertID, err := res.LastInsertId()
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error getting last insert ID", http.StatusInternalServerError)
			return
		}

		order.ID = int(lastInsertID)

		rows, err := tx.Query("SELECT product_id, quantity, price FROM cart_items WHERE cart_id = ?", cartID)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error fetching cart items", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var orderItems []models.OrderItem
		for rows.Next() {
			var orderItem models.OrderItem
			err := rows.Scan(&orderItem.ProductID, &orderItem.Quantity, &orderItem.Price)
			if err != nil {
				tx.Rollback()
				http.Error(w, "Error scanning cart item", http.StatusInternalServerError)
				return
			}
			orderItem.OrderID = order.ID
			orderItems = append(orderItems, orderItem)
		}

		for _, orderItem := range orderItems {
			_, err = tx.Exec("INSERT INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)", orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price)
			if err != nil {
				tx.Rollback()
				http.Error(w, "Error inserting order item", http.StatusInternalServerError)
				return
			}
		}

		for _, orderItem := range orderItems {
			var existQuantity int
			err = db.DB.QueryRow("SELECT quantity FROM products WHERE id = ?", orderItem.ProductID).Scan(&existQuantity)
			if err != nil {
				tx.Rollback()
				http.Error(w, "Error fetching product quantity", http.StatusInternalServerError)
				return
			}
			if existQuantity < orderItem.Quantity {
				tx.Rollback()
				http.Error(w, "Not enough product quantity", http.StatusBadRequest)
				return
			}
		}

		for _, orderItem := range orderItems {
			_, err = tx.Exec("UPDATE products SET quantity = quantity - ? WHERE id = ?", orderItem.Quantity, orderItem.ProductID)
			if err != nil {
				tx.Rollback()
				http.Error(w, "Error inserting order item", http.StatusInternalServerError)
				return
			}
		}

		_, err = tx.Exec("DELETE FROM cart_items WHERE cart_id = ?", cartID)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error clearing cart", http.StatusInternalServerError)
			return
		}

		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			http.Error(w, "Transaction commit error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(order)
	})
}

// GetOrders godoc
// @Summary Get all orders
// @Description Get all orders for the authenticated user
// @Tags orders
// @Produce  json
// @Success 200 {array} models.Order
// @Failure 500 {string} string "Internal server error"
// @Router /orders [get]
func (db *AppHandler) GetOrders() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)

		var orders []models.Order
		rows, err := db.DB.Query("SELECT id, user_id, total_price, created_at FROM orders WHERE user_id = ?", userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var order models.Order
			if err := rows.Scan(&order.ID, &order.UserID, &order.TotalPrice, &order.CreatedAt); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			orders = append(orders, order)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orders)
	})
}

// GetOrderItems godoc
// @Summary Get all items for a specific order
// @Description Get all items for a specific order
// @Tags orders
// @Produce  json
// @Param order_id path int true "Order ID"
// @Success 200 {array} models.OrderItem
// @Failure 500 {string} string "Internal server error"
// @Router /orders/{order_id} [get]
func (db *AppHandler) GetOrderItems() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orderID := vars["order_id"]

		var orderItems []models.OrderItem

		rows, err := db.DB.Query("SELECT id, order_id, product_id, quantity, price FROM order_items WHERE order_id = ?", orderID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var orderItem models.OrderItem
			if err := rows.Scan(&orderItem.ID, &orderItem.OrderID, &orderItem.ProductID, &orderItem.Quantity, &orderItem.Price); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			orderItems = append(orderItems, orderItem)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orderItems)
	})
}

// UpdateOrderStatus godoc
// @Summary Update the status of an order
// @Description Update the status of an order
// @Tags orders
// @Produce  json
// @Param order_id path int true "Order ID"
// @Param status body string true "Order Status"
// @Success 200 {string} string "Order status updated"
// @Failure 403 {string} string "Only admin can update order status"
// @Failure 500 {string} string "Internal server error"
// @Router /orders/{order_id}/status [put]
func (db *AppHandler) UpdateOrderStatus() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userRole := r.Context().Value("role").(string)
		if userRole != "admin" {
			http.Error(w, "Only admin can update order status", http.StatusForbidden)
			return
		}

		vars := mux.Vars(r)
		orderID := vars["order_id"]

		var req struct {
			Status string `json:"status"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		_, err := db.DB.Exec("UPDATE orders SET status = ? WHERE id = ?", req.Status, orderID)
		if err != nil {
			http.Error(w, "Error updating order status", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Order status updated"})
	})
}
