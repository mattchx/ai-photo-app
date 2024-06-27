package main

import (
	"ai-photo-app/db"
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
	router.Post("/logout", handler.MakeHandler(handler.HandleLogoutCreate))
	router.Post("/login", handler.MakeHandler(handler.HandleLoginCreate))
	router.Get("/auth/callback", handler.MakeHandler(handler.HandleAuthCallback))
	router.Get("/account/setup", handler.MakeHandler(handler.HandleAccountSetupIndex))
	
	router.Post("/account/setup", handler.MakeHandler(handler.HandleAccountSetupCreate))

	router.Group(func(auth chi.Router) {
		auth.Use(handler.WithAccountSetup)
		auth.Get("/settings", handler.MakeHandler(handler.HandleSettingsIndex))
		auth.Put("/settings/account/profile", handler.MakeHandler(handler.HandleSettingsUsernameUpdate))
	})

	port := os.Getenv("HTTP_LISTEN_ADDR")
	slog.Info("App is running on ", "port", port)
	log.Fatal(http.ListenAndServe(os.Getenv("HTTP_LISTEN_ADDR"), router))
}

func initEverthing() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	if err := db.Init(); err != nil {
		return err
	}
	return sb.Init()
}
