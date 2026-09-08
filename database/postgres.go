package database

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

var Db *sqlx.DB

// ConnectDatabase loads .env, builds the connection string, and opens the DB.
func ConnectDatabase() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: could not load .env file:", err)
	}

	host := os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	pass := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, pass,
	)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	Db = db
	fmt.Println("Successfully connected to database!")
}
