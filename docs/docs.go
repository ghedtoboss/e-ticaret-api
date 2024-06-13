// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/admin/orders": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get all orders by admin",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "Get all orders by admin",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Order"
                            }
                        }
                    },
                    "403": {
                        "description": "Only admin can access this endpoint",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Error fetching orders",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/admin/products": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Add a new product by admin",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "Add a product by admin",
                "parameters": [
                    {
                        "description": "Product",
                        "name": "product",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Product"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.Product"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "403": {
                        "description": "Only admin can add products",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Error adding product",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/admin/users": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get all users by admin",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "Get all users",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.User"
                            }
                        }
                    },
                    "403": {
                        "description": "Only admin can access this endpoint",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Error fetching users",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/cart": {
            "get": {
                "description": "Get all items in the cart for the authenticated user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cart"
                ],
                "summary": "Get all items in the cart",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.CartItem"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "Add a product to the cart",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cart"
                ],
                "summary": "Add a product to the cart",
                "parameters": [
                    {
                        "description": "Cart Item",
                        "name": "cartItem",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.CartItem"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.CartItem"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/carts/decrease/{item_id}": {
            "put": {
                "description": "Decrease the quantity of an item in the cart",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cart"
                ],
                "summary": "Decrease the quantity of an item in the cart",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Item ID",
                        "name": "item_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ürün adeti azaldı.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Cart not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/carts/increase/{item_id}": {
            "put": {
                "description": "Increase the quantity of an item in the cart",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cart"
                ],
                "summary": "Increase the quantity of an item in the cart",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Item ID",
                        "name": "item_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ürün adeti arttırıldı.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/carts/remove/cart/items": {
            "delete": {
                "description": "Remove all items from the cart",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cart"
                ],
                "summary": "Remove all items from the cart",
                "responses": {
                    "200": {
                        "description": "Sepet temizlendi.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/carts/remove/{item_id}": {
            "delete": {
                "description": "Remove an item from the cart",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cart"
                ],
                "summary": "Remove an item from the cart",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Item ID",
                        "name": "item_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ürün sepetten kaldırıldı.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Item not found in cart",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/order": {
            "post": {
                "description": "Create an order for the authenticated user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "orders"
                ],
                "summary": "Create an order",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.Order"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Cart not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/orders": {
            "get": {
                "description": "Get all orders for the authenticated user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "orders"
                ],
                "summary": "Get all orders",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Order"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/orders/{order_id}": {
            "get": {
                "description": "Get all items for a specific order",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "orders"
                ],
                "summary": "Get all items for a specific order",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Order ID",
                        "name": "order_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.OrderItem"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/orders/{order_id}/status": {
            "put": {
                "description": "Update the status of an order",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "orders"
                ],
                "summary": "Update the status of an order",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Order ID",
                        "name": "order_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Order Status",
                        "name": "status",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Order status updated",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "403": {
                        "description": "Only admin can update order status",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/product": {
            "post": {
                "description": "Add a new product by seller",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Add a new product",
                "parameters": [
                    {
                        "description": "Product to add",
                        "name": "product",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Product"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.Product"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/product/{id}": {
            "put": {
                "description": "Update an existing product by seller",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Update an existing product",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Product ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Product to update",
                        "name": "product",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Product"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content",
                        "schema": {
                            "$ref": "#/definitions/models.Product"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Product not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a product by seller",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Delete a product",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Product ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Product deleted",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/products": {
            "get": {
                "description": "Get all products with optional filters",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Get all products",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Category",
                        "name": "category",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Search term",
                        "name": "search",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Sort by",
                        "name": "sort_by",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Order (asc or desc)",
                        "name": "order",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Product"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.CartItem": {
            "description": "Sepet öğesi modelini temsil eder",
            "type": "object",
            "properties": {
                "cart_id": {
                    "type": "integer",
                    "example": 1
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "price": {
                    "type": "number",
                    "example": 19.99
                },
                "product_id": {
                    "type": "integer",
                    "example": 1
                },
                "quantity": {
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "models.Order": {
            "description": "Sipariş modelini temsil eder",
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "status": {
                    "type": "string",
                    "example": "pending"
                },
                "total_price": {
                    "type": "number",
                    "example": 199.99
                },
                "user_id": {
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "models.OrderItem": {
            "description": "Sipariş öğesi modelini temsil eder",
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "order_id": {
                    "type": "integer",
                    "example": 1
                },
                "price": {
                    "type": "number",
                    "example": 99.99
                },
                "product_id": {
                    "type": "integer",
                    "example": 1
                },
                "quantity": {
                    "type": "integer",
                    "example": 2
                }
            }
        },
        "models.Product": {
            "description": "Ürün modelini temsil eder",
            "type": "object",
            "properties": {
                "category": {
                    "type": "string",
                    "example": "Electronics"
                },
                "description": {
                    "type": "string",
                    "example": "Description"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "image_url": {
                    "type": "string",
                    "example": "http://..."
                },
                "name": {
                    "type": "string",
                    "example": "Product Name"
                },
                "price": {
                    "type": "number",
                    "example": 19.99
                },
                "quantity": {
                    "type": "integer",
                    "example": 100
                },
                "seller_id": {
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "models.User": {
            "description": "Kullanıcı modelini temsil eder",
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "example": "user@example.com"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "name": {
                    "type": "string",
                    "example": "John Doe"
                },
                "password": {
                    "type": "string"
                },
                "role": {
                    "type": "string",
                    "example": "seller"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "E-Ticaret API",
	Description:      "E-Ticaret API dokümantasyonu.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
