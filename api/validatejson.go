package api

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
)

// ValidatePassword checks if the password meets the specified criteria.
var ValidatePassword validator.Func = func(fl validator.FieldLevel) bool {
	password := fl.Field().Interface().(string)

	// Check if the password is at least 8 characters long
	if utf8.RuneCountInString(password) < 8 {
		return false
	}

	// Check if the password contains at least one digit and one symbol
	hasDigit := false
	hasSymbol := false
	hasUpper := false
	hasLower := false
	for _, char := range password {
		if unicode.IsNumber(char) {
			hasDigit = true
		}
		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSymbol = true
		}
		if unicode.IsUpper(char) {
			hasUpper = true
		}
		if unicode.IsLower(char) {
			hasLower = true
		}
	}

	//fmt.Println("Validating password:", password)

	return hasDigit && hasSymbol && hasUpper && hasLower
}

// ImageURLValidation is a custom validator function to check if the URL points to an image.
var ImageURLValidation validator.Func = func(fl validator.FieldLevel) bool {
	urlStr := fl.Field().Interface().(string)

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
	//defer resp.Body.Close()

	// Check if the content type indicates an image
	contentType := resp.Header.Get("Content-Type")

	fmt.Println(contentType)

	err = resp.Body.Close()
	if err != nil {
		log.Fatal(err.Error())
	}

	return strings.HasPrefix(contentType, "image/")

}
