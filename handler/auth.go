package handler

import (
	"ai-photo-app/pkg/sb"
	"ai-photo-app/pkg/util"
	"ai-photo-app/view/auth"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

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

	setAuthCookie(w, resp.AccessToken)
	// http.Redirect(w, r, "/", http.StatusSeeOther)
	return hxRedirect(w, r, "/")
}

func validateSignup(params auth.SignupParams) (bool, auth.SignupErrors) {

	okay := true
	errors := auth.SignupErrors{}

	// check if email is valid
	if ok := util.IsValidEmail(params.Email); !ok {
		okay = false
		errors.Email = "Invalid email address"
	}

	// check if password is valid
	if str, ok := util.ValidatePassword(params.Password); !ok {
		okay = false
		errors.Password = str
	}
	// check if confirm password is valid
	if str, ok := util.ValidatePassword(params.ConfirmPassword); !ok {
		okay = false
		errors.ConfirmPassword = str
	}
	// check if passwords match
	if strings.Compare(params.Password, params.ConfirmPassword) != 0 {
		// if params.Password == params.ConfirmPassword {
		okay = false
		errors.ConfirmPassword = "Passwords do not match"
	}

	fmt.Print(okay, errors)
	return okay, errors
}

func HandleSignupCreate(w http.ResponseWriter, r *http.Request) error {
	params := auth.SignupParams{
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirmPassword"),
	}
	okay, errors := validateSignup(params)
	if !okay {
		return render(r, w, auth.SignupForm(params, errors))
	}
	user, err := sb.Client.Auth.SignUp(r.Context(), supabase.UserCredentials{
		Email:    params.Email,
		Password: params.Password,
	})
	if err != nil {
		fmt.Println("There was an error: ", err)
		return err
	}

	return render(r, w, auth.SignupSuccess(user.Email))
}

func HandleAuthCallback(w http.ResponseWriter, r *http.Request) error {

	accessToken := r.URL.Query().Get("access_token")
	if len(accessToken) == 0 {
		return render(r, w, auth.CallbackScript())
	}
	setAuthCookie(w, accessToken)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func HandleLogoutCreate(w http.ResponseWriter, r *http.Request) error {
	cookie := http.Cookie{
		Value:    "",
		Name:     "at",
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return nil
}

func setAuthCookie(w http.ResponseWriter, accessToken string) error {
	cookie := &http.Cookie{
		Value:    accessToken,
		Name:     "at",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)

	return nil
}
