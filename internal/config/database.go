package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		panic("DATABASE_URL is empty")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(fmt.Sprintf("failed to load DB driver: %v", err))
	}

	if err := db.Ping(); err != nil {
		panic(fmt.Sprintf("database is not responding: %v", err))
	}

	fmt.Println("success connect to database")
	return db
}
