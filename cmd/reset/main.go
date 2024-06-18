package main

import (
	"database/sql"
	"ai-photo-app/db"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func createDB() (*sql.DB, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	var (
		host   = os.Getenv("DB_HOST")
		user   = os.Getenv("DB_USER")
		pass   = os.Getenv("DB_PASSWORD")
		dbname = os.Getenv("DB_NAME")
	)
	return db.CreateDatabase(dbname, user, pass, host)
}

func main() {
	db, err := createDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tables := []string{
		"schema_migrations",
		"accounts",
	}

	for _, table := range tables {
		query := fmt.Sprintf("drop table if exists %s cascade", table)
		if _, err := db.Exec(query); err != nil {
			log.Fatal(err)
		}
	}
}
