package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
)

// ValidatePassword checks if the password meets the specified criteria.
func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Check if the password is at least 8 characters long
	if utf8.RuneCountInString(password) < 8 {
		return false
	}

	// Check if the password contains at least one digit and one symbol
	hasDigit := false
	hasSymbol := false
	hasUpper := false
	for _, char := range password {
		if unicode.IsDigit(char) && (unicode.IsPunct(char) || unicode.IsSymbol(char)) && unicode.IsUpper(char) {
			hasDigit = true
			hasSymbol = true
			hasUpper = true
		}
	}

	fmt.Println("Validating password:", password)

	return hasDigit && hasSymbol && hasUpper
}

// ImageURLValidation is a custom validator function to check if the URL points to an image.
func ImageURLValidation(fl validator.FieldLevel) bool {
	urlStr := fl.Field().String()

	// Parse the URL
	u, err := url.Parse(urlStr)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return isImageURL(u)
}

func isImageURL(u *url.URL) bool {

	resp, err := http.Get(u.String())
	if err != nil {
		return false
	}
	fmt.Println(resp)
	defer resp.Body.Close()

	// Check if the content type indicates an image
	contentType := resp.Header.Get("Content-Type")

	fmt.Println(contentType)

	return strings.HasPrefix(contentType, "image/")

}

// Register the custom validation function
var V *validator.Validate
