package handler

import "net/http"
import "ai-photo-app/view/settings"

func HandleSettingsIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	return render(r, w, settings.Index(user))
}