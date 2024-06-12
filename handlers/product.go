package handlers

import (
	"e-ticaret-api/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (db *AppHandler) AddProduct() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		UserID := r.Context().Value("userID").(int)
		UserRole := r.Context().Value("role").(string)
		if UserRole != "seller" {
			http.Error(w, "Sadece satıcılar ürün ekleyebilir.", http.StatusUnauthorized)
			return
		}

		var product models.Product
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, err := db.DB.Exec("INSERT INTO products (name, description, quantity, price, seller_id, category, image_url) VALUES (?, ?, ?, ?, ?, ?, ?)", product.Name, product.Description, product.Quantity, product.Price, UserID, product.Category, product.ImageURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		lastInsertID, err := res.LastInsertId()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		product.ID = int(lastInsertID)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(product)

	})
}

func (db *AppHandler) UpdateProduct() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		productID := vars["id"]

		var product models.Product
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var existProduct models.Product
		row := db.DB.QueryRow("SELECT name, description, quantity, price, seller_id, image_url FROM products WHERE id = ?", productID)
		if err := row.Scan(&existProduct.Name, &existProduct.Description, &existProduct.Quantity, &existProduct.Price, &existProduct.SellerID, &existProduct.ImageURL); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if product.Name != "" {
			existProduct.Name = product.Name
		}
		if product.Description != "" {
			existProduct.Description = product.Description
		}
		if product.Quantity != existProduct.Quantity {
			existProduct.Quantity = product.Quantity
		}
		if product.Price != existProduct.Price {
			existProduct.Price = product.Price
		}
		if product.ImageURL != "" {
			existProduct.ImageURL = product.ImageURL
		}

		_, err := db.DB.Exec("UPDATE products SET name = ?, description = ?, quantity = ?, price = ?, image_url = ? WHERE id = ?", existProduct.Name, existProduct.Description, existProduct.Quantity, existProduct.Price, existProduct.ImageURL, productID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(existProduct)

	})
}

func (db *AppHandler) DeleteProduct() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		productID := vars["id"]

		UserRole := r.Context().Value("role").(string)
		if UserRole != "seller" {
			http.Error(w, "Sadece satıcılar ürün silinebilir.", http.StatusUnauthorized)
			return
		}

		_, err := db.DB.Exec("DELETE FROM products WHERE id = ?", productID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Urun silindi."})
	})
}

func (db *AppHandler) GetProducts() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		category := query.Get("category")
		search := query.Get("search")
		sortBy := query.Get("sort_by")
		order := query.Get("order")

		baseQuery := "SELECT id, name, description, quantity, price, seller_id, category, image_url FROM products WHERE 1=1"

		if category != "" {
			baseQuery += " AND category = ?"
		}
		if search != "" {
			baseQuery += " AND (name LIKE ? OR description LIKE ?)"
			search = "%" + search + "%"
		}
		if sortBy != "" {
			baseQuery += " ORDER BY " + sortBy
			if order != "" {
				baseQuery += " " + order
			} else {
				baseQuery += " ASC"
			}
		}

		stmt, err := db.DB.Prepare(baseQuery)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		rows, err := stmt.Query(category, search, search)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var products []models.Product
		for rows.Next() {
			var product models.Product
			if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Quantity, &product.Price, &product.SellerID, &product.Category, &product.ImageURL); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			products = append(products, product)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	})
}
