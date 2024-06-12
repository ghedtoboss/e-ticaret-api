package handlers

import (
	"e-ticaret-api/models"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func (db *AppHandler) Register() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Register handler called")
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println("JSON Decode error: ", err)
			return
		}
		log.Println("User data decoded: ", user)

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Error hashing password: ", err)
			return
		}
		user.Password = string(hashedPassword)
		log.Println("Password hashed:", user.Password)

		_, err = db.DB.Exec("INSERT INTO users (email, password, name, role) VALUES (?, ?, ?, ?)", user.Email, user.Password, user.Name, user.Role)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Database Insert Error: ", err)
			return
		}
		log.Println("User inserted into database: ", user)

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(user); err != nil {
			log.Println("Error encoding response: ", err)
		}
	})
}

func (db *AppHandler) Login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var creds models.User
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println("JSON Decode error: ", err)
			return
		}

		var storedUser models.User
		row := db.DB.QueryRow("SELECT id, email, password, name, role FROM users WHERE email = ?", creds.Email)
		if err := row.Scan(&storedUser.ID, &storedUser.Email, &storedUser.Password, &storedUser.Name, &storedUser.Role); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(creds.Password)); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			log.Println("Password comparison error: ", err)
			return
		}

		expirationTime := time.Now().Add(24 * time.Hour)
		claims := &models.Claims{
			Username: storedUser.Email,
			UserID:   storedUser.ID,
			Role:     storedUser.Role,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Error signing token: ", err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})

		if err := json.NewEncoder(w).Encode(map[string]string{"token": tokenString}); err != nil {
			log.Println("Error encoding response: ", err)
		}
	})
}
