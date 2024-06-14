# E-Commerce API

This is an E-Commerce API built with Go, providing endpoints for user authentication, product management, cart management, order management, returns, and reviews.

## Features

1. User Registration and Authentication
2. Product Management (Add, Update, Delete)
3. Cart Management (Add to Cart, Remove from Cart, Increase/Decrease Item Quantity)
4. Order Management (Create Order, View Orders, View Order Items)
5. Return Management (Create Return, View Returns)
6. Review Management (Create Review, View Reviews)
7. Admin Management (View Users, Add Products, View All Orders)

## Technologies Used

- Go (Golang)
- Gorilla Mux
- MySQL
- JWT for Authentication
- Swag for API Documentation

## Installation

1. Clone the repository:

```sh
git clone https://github.com/yourusername/e-commerce-api.git
cd e-commerce-api

2. Create a '.env' file and add your environment variables:
DATABASE_URL="your_database_url"
JWT_SECRET_KEY="your_jwt_secret_key"

3. Install the dependencies:
go mod tidy

4. Generate Swagger documentation:
swag init

5. Run the application:
go run main.go


API Endpoints
Authentication
POST /register: Register a new user
POST /login: Login and get a JWT token
Products
POST /product: Add a new product (Seller only)
PUT /product/{id}: Update a product (Seller only)
DELETE /product/{id}: Delete a product (Seller only)
GET /products: Get a list of products
Cart
POST /cart: Add an item to the cart
GET /cart: Get cart items
DELETE /carts/remove/{item_id}: Remove an item from the cart
PUT /carts/decrease/{item_id}: Decrease item quantity in the cart
PUT /carts/increase/{item_id}: Increase item quantity in the cart
DELETE /carts/remove/cart/items: Clear all items in the cart
Orders
POST /order: Create a new order
GET /orders: Get user orders
GET /orders/{order_id}: Get items of a specific order
PUT /orders/{order_id}/status: Update the status of an order (Admin only)
Returns
POST /returns: Create a return
GET /returns: Get returns
Reviews
POST /reviews: Create a review
GET /reviews/{product_id}: Get reviews for a product
Admin
GET /admin/users: Get all users (Admin only)
POST /admin/products: Add a product (Admin only)
GET /admin/orders: Get all orders (Admin only)
Swagger Documentation
The API documentation can be accessed at /swagger/index.html after running the application.
