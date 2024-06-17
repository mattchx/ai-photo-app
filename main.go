package main

import (
	"ai-photo-app/handler"
	"ai-photo-app/pkg/sb"
	"embed"
	"log"
	"log/slog"
	"net/http"
	"os"

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
	router.Use(handler.WithUser)

	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(FS))))

	router.Get("/", handler.MakeHandler(handler.HandleHomeIndex))
	router.Get("/login", handler.MakeHandler(handler.HandleLoginIndex))
	router.Get("/login/provider/google", handler.MakeHandler(handler.HandleLoginWithGoogle))
	router.Get("/signup", handler.MakeHandler(handler.HandleSignupIndex))

	router.Post("/login", handler.MakeHandler(handler.HandleLoginCreate))
	router.Post("/signup", handler.MakeHandler(handler.HandleSignupCreate))
	router.Post("/logout", handler.MakeHandler(handler.HandleLogoutCreate))

	router.Get("/auth/callback", handler.MakeHandler(handler.HandleAuthCallback))

	router.Group(func(auth chi.Router) {
		auth.Use(handler.WithAuth)
		auth.Get("/settings", handler.MakeHandler(handler.HandleSettingsIndex))
	})

	port := os.Getenv("HTTP_LISTEN_ADDR")
	slog.Info("App is running on ", "port", port)
	log.Fatal(http.ListenAndServe(os.Getenv("HTTP_LISTEN_ADDR"), router))
}

func initEverthing() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	return sb.Init()
}
