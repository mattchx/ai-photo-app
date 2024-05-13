package util

import (
	"regexp"
	"unicode"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,8}$`)

func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)

}


// validatePassword checks if the password is strong and meets the criteria:
// - At least 8 characters long
// - Contains at least one digit
// - Contains at least one lowercase letter
// - Contains at least one uppercase letter
// - Contains at least one special character
func ValidatePassword(password string) (string, bool) {
	var (
		hasUpper     = false
		hasLower     = false
		hasNumber    = false
		hasSpecial   = false
		specialRunes = "!@#$%^&*"
	)

	if len(password) < 8 {
		return "Password must contain at least 8 characters", false
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char) || strings.ContainsRune(specialRunes, char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return "Password must contain at least 1 uppercase character", false
	}
	if !hasLower {
		return "Password must contain at least 1 lowercase character", false
	}
	if !hasNumber {
		return "Password must contain at least 1 numeric character (0, 1, 2, ...)", false
	}
	if !hasSpecial {
		return "Password must contain at least 1 special character (@, ;, _, ...)", false
	}
	return "", true
}

