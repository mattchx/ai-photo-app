package handler

import (
	"ai-photo-app/pkg/sb"
	"ai-photo-app/pkg/util"
	"ai-photo-app/view/auth"
	"log/slog"
	"net/http"

	"github.com/nedpals/supabase-go"
)

func HandleLoginIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.Login())
}

func HandleLoginCreate(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	if !util.IsValidEmail(credentials.Email) {
		return render(r, w, auth.LoginForm(credentials, auth.LoginErrors{
			Email: "Please enter a valid email",
		}))
	}

	if reason, ok := util.ValidatePassword(credentials.Password); !ok {
		return render(r, w, auth.LoginForm(credentials, auth.LoginErrors{
			Password: reason,
		}))
	}

	resp, error := sb.Client.Auth.SignIn(r.Context(), credentials)
	if error != nil {
		slog.Error("login error", "error", error)
		return render(r, w, auth.LoginForm(credentials, auth.LoginErrors{
			InvalidCredentials: "The credentials entered are invalid",
		}))
	}

	cookie := &http.Cookie{
		Value: resp.AccessToken,
		Name:  "at",
		Path:  "/",
		HttpOnly: true,
		Secure: true,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}
