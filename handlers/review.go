// handlers/review.go
package handlers

import (
	"e-ticaret-api/models"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// CreateReview godoc
// @Summary Create a product review
// @Description Create a new review for a product
// @Tags Reviews
// @Accept  json
// @Produce  json
// @Param review body models.Review true "Review"
// @Success 201 {object} models.Review
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /reviews [post]
func (db *AppHandler) CreateReview() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)

		var review models.Review
		if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		review.UserID = userID
		review.CreatedAt = time.Now()

		_, err := db.DB.Exec("INSERT INTO reviews (product_id, user_id, rating, comment, created_at) VALUES (?, ?, ?, ?, ?)",
			review.ProductID, review.UserID, review.Rating, review.Comment, review.CreatedAt)
		if err != nil {
			http.Error(w, "Error creating review", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(review)
	})
}

// GetReviews godoc
// @Summary Get product reviews
// @Description Get a list of reviews for a product
// @Tags Reviews
// @Produce  json
// @Param product_id path int true "Product ID"
// @Success 200 {array} models.Review
// @Failure 500 {string} string "Internal server error"
// @Router /reviews/{product_id} [get]
func (db *AppHandler) GetReviews() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		productID := vars["product_id"]

		rows, err := db.DB.Query("SELECT id, product_id, user_id, rating, comment, created_at FROM reviews WHERE product_id = ?", productID)
		if err != nil {
			http.Error(w, "Error fetching reviews", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var reviews []models.Review
		for rows.Next() {
			var review models.Review
			if err := rows.Scan(&review.ID, &review.ProductID, &review.UserID, &review.Rating, &review.Comment, &review.CreatedAt); err != nil {
				http.Error(w, "Error scanning review", http.StatusInternalServerError)
				return
			}
			reviews = append(reviews, review)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reviews)
	})
}
