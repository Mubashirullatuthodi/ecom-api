package utils

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

func FirstLetterCapital(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if len(value) == 0 {
		return true
	}
	return unicode.IsUpper(rune(value[0]))
}

func ValidatePassword(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	var hasMinLen, hasUpper, hasSpecial bool
	var specialCharRE = regexp.MustCompile(`[!@#$%^&*()_+|]{1}`)

	if len(value) >= 8 {
		hasMinLen = true
	}
	for _, char := range value {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case specialCharRE.MatchString(string(char)):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasSpecial
}

type NewUser struct {
	FirstName string `json:"firstName" validate:"required,firstLetterCapital"`
	LastName  string `json:"lastName"`
	Email     string `json:"email" validate:"required,email"`
	Gender    string `json:"gender"`
	Phone     string `json:"phoneNo" validate:"required,len=10,numeric"`
	Password  string `json:"password" validate:"required,password"`
	Status    string `json:"status"`
}

// validate the userInstance
func ValidateUserInstance(user NewUser) map[string]string {
	validate := validator.New()

	// Custom validation
	validate.RegisterValidation("firstLetterCapital", FirstLetterCapital)
	validate.RegisterValidation("password", ValidatePassword)

	// Validate user instance
	err := validate.Struct(user)
	if err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Password":
				if err.Tag() == "password" {
					errors["Password"] = "Password must be at least 8 characters long, include at least one uppercase letter and one special character."
				} else {
					errors[err.Field()] = err.Tag()
				}
			case "Email":
				if err.Tag() == "email" {
					errors["Email"] = "It must be a valid email address."
				} else {
					errors[err.Field()] = err.Tag()
				}
			case "FirstName":
				if err.Tag() == "firstLetterCapital" {
					errors["FirstName"] = "First letter must be capital."
				} else {
					errors[err.Field()] = err.Tag()
				}
			default:
				errors[err.Field()] = err.Error()
			}
		}
		return errors
	}

	return nil
}
