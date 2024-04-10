package main

import (
	"ai-photo-app/handler"
	"log"
	"log/slog"
	"net/http"
	"os"
	"embed"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

//go:embed public
var FS embed.FS

func main() {

	if err := initEverthing(); err != nil {
		log.Fatal(err)
	}
	router := chi.NewMux()

	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(FS))))
	router.Get("/", handler.MakeHandler(handler.HandleHomeIndex))

	port := os.Getenv("HTTP_LISTEN_ADDR")
	slog.Info("App is running on ", "port", port)
	log.Fatal(http.ListenAndServe(os.Getenv("HTTP_LISTEN_ADDR"), router))
}

func initEverthing() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	return godotenv.Load()
}