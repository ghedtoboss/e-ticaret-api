package handlers

import (
	"database/sql"
)

type AppHandler struct {
	DB *sql.DB
}
