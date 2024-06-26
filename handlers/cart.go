package handlers

import (
	"e-ticaret-api/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// AddToCart godoc
// @Summary Add a product to the cart
// @Description Add a product to the cart
// @Tags cart
// @Accept  json
// @Produce  json
// @Param cartItem body models.CartItem true "Cart Item"
// @Success 201 {object} models.CartItem
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /cart [post]
func (db *AppHandler) AddToCart() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)

		var CartItem models.CartItem
		if err := json.NewDecoder(r.Body).Decode(&CartItem); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//Kullanıcının sepeti var mı
		var cartID int
		err := db.DB.QueryRow("SELECT id FROM carts WHERE user_id = ?", userID).Scan(&cartID)
		if err != nil {
			//sepet yoksa sepet oluşturma
			res, err := db.DB.Exec("INSERT INTO carts (user_id) VALUES (?)", userID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			cartID64, _ := res.LastInsertId()
			cartID = int(cartID64)
		}

		CartItem.CartID = cartID
		_, err = db.DB.Exec("INSERT INTO cart_items (cart_id, product_id, quantity, price) VALUES (?, ?, ?, ?)",
			CartItem.CartID, CartItem.ProductID, CartItem.Quantity, CartItem.Price)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(CartItem)
	})
}

// GetCartItems godoc
// @Summary Get all items in the cart
// @Description Get all items in the cart for the authenticated user
// @Tags cart
// @Produce  json
// @Success 200 {array} models.CartItem
// @Failure 500 {string} string "Internal server error"
// @Router /cart [get]
func (db *AppHandler) GetCartItems() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)

		var cartID int
		err := db.DB.QueryRow("SELECT id FROM carts WHERE user_id = ?", userID).Scan(&cartID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rows, err := db.DB.Query("SELECT id, cart_id, product_id, quantity, price FROM cart_items WHERE cart_id = ?", cartID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var cartItems []models.CartItem
		for rows.Next() {
			var cartItem models.CartItem
			if err := rows.Scan(&cartItem.ID, &cartItem.CartID, &cartItem.ProductID, &cartItem.Quantity, &cartItem.Price); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			cartItems = append(cartItems, cartItem)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cartItems)
	})
}

// RemoveFromCart godoc
// @Summary Remove an item from the cart
// @Description Remove an item from the cart
// @Tags cart
// @Produce  json
// @Param item_id path int true "Item ID"
// @Success 200 {string} string "Ürün sepetten kaldırıldı."
// @Failure 404 {string} string "Item not found in cart"
// @Failure 500 {string} string "Internal server error"
// @Router /carts/remove/{item_id} [delete]
func (db *AppHandler) RemoveFromCart() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		itemID := vars["item_id"]

		userID := r.Context().Value("userID").(int)
		var cartID int
		err := db.DB.QueryRow("SELECT id FROM carts WHERE user_id = ?", userID).Scan(&cartID)
		if err != nil {
			http.Error(w, "Cart not found", http.StatusNotFound)
			return
		}

		// İlk olarak, itemID'nin doğru cart_id ile eşleşip eşleşmediğini kontrol edelim
		var dbCartID int
		err = db.DB.QueryRow("SELECT cart_id FROM cart_items WHERE id = ? AND cart_id = ?", itemID, cartID).Scan(&dbCartID)
		if err != nil || dbCartID != cartID {
			http.Error(w, "Item not found in cart", http.StatusNotFound)
			return
		}

		// Eğer itemID ve cartID eşleşiyorsa, ürünü sil
		_, err = db.DB.Exec("DELETE FROM cart_items WHERE id = ? AND cart_id = ?", itemID, cartID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Ürün sepetten kaldırıldı."})
	})
}

// DecreaseItemQuantity godoc
// @Summary Decrease the quantity of an item in the cart
// @Description Decrease the quantity of an item in the cart
// @Tags cart
// @Produce  json
// @Param item_id path int true "Item ID"
// @Success 200 {string} string "Ürün adeti azaldı."
// @Failure 404 {string} string "Cart not found"
// @Failure 500 {string} string "Internal server error"
// @Router /carts/decrease/{item_id} [put]
func (db *AppHandler) DecreaseItemQuantity() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		itemID := vars["item_id"]

		userID := r.Context().Value("userID").(int)

		var cartID int
		err := db.DB.QueryRow("SELECT id FROM carts WHERE user_id = ?", userID).Scan(&cartID)
		if err != nil {
			http.Error(w, "Cart not found", http.StatusInternalServerError)
			return
		}

		tx, err := db.DB.Begin()
		if err != nil {
			http.Error(w, "Transaction begin error", http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec("UPDATE cart_items SET quantity = quantity - 1 WHERE id = ? AND cart_id = ?", itemID, cartID)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error decreasing quantity", http.StatusInternalServerError)
			return
		}

		var quantity, productID int
		err = tx.QueryRow("SELECT quantity, product_id FROM cart_items WHERE id = ? AND cart_id = ?", itemID, cartID).Scan(&quantity, &productID)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error fetching quantity", http.StatusInternalServerError)
			return
		}
		log.Println("Quantity:", quantity, "Product ID:", productID)

		if quantity == 0 {
			_, err = tx.Exec("DELETE FROM cart_items WHERE id = ? AND cart_id = ?", itemID, cartID)
			if err != nil {
				tx.Rollback()
				http.Error(w, "Error deleting item", http.StatusInternalServerError)
				return
			}
		} else {
			var productPrice float64
			err = tx.QueryRow("SELECT price FROM products WHERE id = ?", productID).Scan(&productPrice)
			if err != nil {
				tx.Rollback()
				http.Error(w, "Error fetching product price", http.StatusInternalServerError)
				return
			}
			log.Println("Product Price:", productPrice)

			_, err = tx.Exec("UPDATE cart_items SET price = ? * quantity WHERE id = ? AND cart_id = ?", productPrice, itemID, cartID)
			if err != nil {
				tx.Rollback()
				http.Error(w, "Error updating price", http.StatusInternalServerError)
				return
			}
		}

		err = tx.Commit()
		if err != nil {
			http.Error(w, "Transaction commit error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if quantity == 0 {
			json.NewEncoder(w).Encode(map[string]string{"message": "Ürün sepetten kaldırıldı."})
		} else {
			json.NewEncoder(w).Encode(map[string]string{"message": "Ürün adeti azaldı."})
		}
	})
}

// IncreaseItemQuantity godoc
// @Summary Increase the quantity of an item in the cart
// @Description Increase the quantity of an item in the cart
// @Tags cart
// @Produce  json
// @Param item_id path int true "Item ID"
// @Success 200 {string} string "Ürün adeti arttırıldı."
// @Failure 500 {string} string "Internal server error"
// @Router /carts/increase/{item_id} [put]
func (db *AppHandler) IncreaseItemQuantity() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		itemID := vars["item_id"]

		userID := r.Context().Value("userID").(int)

		var cartID int
		err := db.DB.QueryRow("SELECT id FROM carts WHERE user_id = ?", userID).Scan(&cartID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tx, err := db.DB.Begin()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec("UPDATE cart_items SET quantity = quantity + 1 WHERE id = ? AND cart_id = ?", itemID, cartID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec("UPDATE cart_items SET price = price * quantity WHERE id = ? AND cart_id = ?", itemID, cartID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tx.Commit()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Ürün adeti arttırıldı."})
	})
}

// RemoveCartItems godoc
// @Summary Remove all items from the cart
// @Description Remove all items from the cart
// @Tags cart
// @Produce  json
// @Success 200 {string} string "Sepet temizlendi."
// @Failure 500 {string} string "Internal server error"
// @Router /carts/remove/cart/items [delete]
func (db *AppHandler) RemoveCartItems() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)

		var cartID int
		err := db.DB.QueryRow("SELECT id FROM carts WHERE user_id = ?", userID).Scan(&cartID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = db.DB.Exec("DELETE FROM cart_items WHERE cart_id = ?", cartID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Sepet temizlendi."})
	})
}
