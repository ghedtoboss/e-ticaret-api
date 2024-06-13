package handlers

import (
	"e-ticaret-api/models"
	"encoding/json"
	"net/http"
	"time"
)

// CreateReturn godoc
// @Summary Create a return request
// @Description Create a new return request for a product
// @Tags Returns
// @Accept  json
// @Produce  json
// @Param return body models.Return true "Return request"
// @Success 201 {object} models.Return
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /returns [post]
func (db *AppHandler) CreateReturn() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var returnReq models.Return
		if err := json.NewDecoder(r.Body).Decode(&returnReq); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		returnReq.Status = "pending"
		returnReq.CreatedAt = time.Now()

		tx, err := db.DB.Begin()
		if err != nil {
			http.Error(w, "Transaction begin error", http.StatusInternalServerError)
			return
		}

		res, err := tx.Exec("INSERT INTO returns (order_id, product_id, reason, status, created_at) VALUES (?, ?, ?, ?, ?)",
			returnReq.OrderID, returnReq.ProductID, returnReq.Reason, returnReq.Status, returnReq.CreatedAt)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error creating return request", http.StatusInternalServerError)
			return
		}

		lastInsertID, err := res.LastInsertId()
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error getting last insert ID", http.StatusInternalServerError)
			return
		}

		returnReq.ID = int(lastInsertID)

		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			http.Error(w, "Transaction commit error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(returnReq)
	})
}

// GetReturns godoc
// @Summary Get return requests
// @Description Get a list of return requests for the current user
// @Tags Returns
// @Produce  json
// @Success 200 {array} models.Return
// @Failure 500 {string} string "Internal server error"
// @Router /returns [get]
func (db *AppHandler) GetReturns() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)

		rows, err := db.DB.Query("SELECT id, order_id, product_id, reason, status, created_at FROM returns WHERE user_id = ?", userID)
		if err != nil {
			http.Error(w, "Error fetching return requests", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var returns []models.Return
		for rows.Next() {
			var returnReq models.Return
			if err := rows.Scan(&returnReq.ID, &returnReq.OrderID, &returnReq.ProductID, &returnReq.Reason, &returnReq.Status, &returnReq.CreatedAt); err != nil {
				http.Error(w, "Error scanning return request", http.StatusInternalServerError)
				return
			}
			returns = append(returns, returnReq)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(returns)
	})
}
