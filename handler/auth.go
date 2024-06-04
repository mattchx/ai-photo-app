package handler

import (
	"ai-photo-app/pkg/kit/validate"
	"ai-photo-app/pkg/sb"
	"ai-photo-app/view/auth"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/nedpals/supabase-go"
)

func HandleLoginIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.Login())
}

func HandleSignupIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.Signup())
}

func HandleLoginCreate(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	resp, error := sb.Client.Auth.SignIn(r.Context(), credentials)
	if error != nil {
		slog.Error("login error", "error", error)
		return render(r, w, auth.LoginForm(credentials, auth.LoginErrors{
			InvalidCredentials: "The credentials entered are invalid",
		}))
	}

	cookie := &http.Cookie{
		Value:    resp.AccessToken,
		Name:     "at",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func HandleSignupCreate(w http.ResponseWriter, r *http.Request) error {
	fmt.Printf("gets here 1")
	params := auth.SignupParams{
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("password-repeat"),
	}
	errors := auth.SignupErrors{}
	if ok := validate.New(&params, validate.Fields{
		"Email": validate.Rules(validate.Email),
		"Password": validate.Rules(validate.Password),
		 "ConfirmPassword": validate.Rules(
				validate.Equal(params.Password)),
		 }).Validate(&errors); !ok {
			slog.Error("signup error", "error", errors)
			return render(r, w, auth.SignupForm(params, errors))
	}
	fmt.Printf("gets here 2")
	user, err := sb.Client.Auth.SignUp(r.Context(), supabase.UserCredentials{
		Email:    params.Email,
		Password: params.Password,
	})
	if err != nil {
		return err
	}
	fmt.Printf("gets here 3")
	return render(r, w, auth.SignupSuccess(user.Email))
}
