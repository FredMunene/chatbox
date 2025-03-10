package utils

import (
	"database/sql"

	"forum/backend/DB"
)

var Database *sql.DB

func Init() {
	Database = DB.CreateConnection()
}
