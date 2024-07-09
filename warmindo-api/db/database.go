package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST") // Assuming you're also using an environment variable for host
	dbname := os.Getenv("POSTGRES_DB")

	connStr := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=require", user, password, host, dbname)

	return sql.Open("postgres", connStr)
}
