package database

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	"github.com/jmoiron/sqlx"
)

var Db *sqlx.DB

// ConnectDatabase loads .env, builds the connection string, and opens the DB.
func ConnectDatabase() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Warning: could not load .env file:", err)
	}

	dsn := viper.GetString("DATABASE_URL")

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

func CreateTasksTable() {
	query := `
		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			completed BOOLEAN NOT NULL DEFAULT FALSE
		);
	`

	if _, err := Db.Exec(query); err != nil {
		panic(err)
	}
}
