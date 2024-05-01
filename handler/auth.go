package handler

import (
	"ai-photo-app/view/auth"
	"net/http"
)

func HandleLogInIndex (w http.ResponseWriter, r *http.Request) error {
	return auth.LogIn().Render(r.Context(), w)
}