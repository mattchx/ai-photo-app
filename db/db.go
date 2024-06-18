package db


import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

var Bun *bun.DB

func CreateDatabase(
	dbname string,
	dbuser string,
	dbpassword string,
	dbhost string,
) (*sql.DB, error) {
	hostArr := strings.Split(dbhost, ":")
	host := hostArr[0]
	port := "5432"
	if len(hostArr) > 1 {
		port = hostArr[1]
	}
	uri := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		dbuser,
		dbpassword,
		dbname,
		host,
		port,
	)
	db, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Init() error {
	var (
		host   = os.Getenv("DB_HOST")
		user   = os.Getenv("DB_USER")
		pass   = os.Getenv("DB_PASSWORD")
		dbname = os.Getenv("DB_NAME")
	)
	db, err := CreateDatabase(dbname, user, pass, host)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	Bun = bun.NewDB(db, pgdialect.New())
	if len(os.Getenv("APP_DEBUG")) > 0 {
		Bun.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}
	return nil
}