package validation

import (
	"regexp"
	"strings"
	"unicode"
)

// Email validation regex
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// ValidateEmail checks if the email format is valid
func ValidateEmail(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	return emailRegex.MatchString(email)
}

// ValidatePassword checks if the password meets security requirements
func ValidatePassword(password string) (bool, []string) {
	var errors []string

	// Check minimum length
	if len(password) < 8 {
		errors = append(errors, "Password must be at least 8 characters long")
	}

	// Check maximum length
	if len(password) > 128 {
		errors = append(errors, "Password must not exceed 128 characters")
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		errors = append(errors, "Password must contain at least one uppercase letter")
	}
	if !hasLower {
		errors = append(errors, "Password must contain at least one lowercase letter")
	}
	if !hasNumber {
		errors = append(errors, "Password must contain at least one number")
	}
	if !hasSpecial {
		errors = append(errors, "Password must contain at least one special character")
	}

	return len(errors) == 0, errors
}

// ValidateUsername checks if the username is valid
func ValidateUsername(username string) (bool, []string) {
	var errors []string

	// Check length
	if len(username) < 3 {
		errors = append(errors, "Username must be at least 3 characters long")
	}
	if len(username) > 30 {
		errors = append(errors, "Username must not exceed 30 characters")
	}

	// Check allowed characters (alphanumeric, underscore, hyphen)
	allowedChars := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !allowedChars.MatchString(username) {
		errors = append(errors, "Username can only contain letters, numbers, underscores, and hyphens")
	}

	// Check if it starts with a letter
	if len(username) > 0 && !unicode.IsLetter(rune(username[0])) {
		errors = append(errors, "Username must start with a letter")
	}

	return len(errors) == 0, errors
}

// ValidatePhoneNumber checks if the phone number format is valid
func ValidatePhoneNumber(phone string) bool {
	// Remove spaces and special characters
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")

	// Check if it starts with + for international format
	phone = strings.TrimPrefix(phone, "+")

	// Check if all remaining characters are digits
	phoneRegex := regexp.MustCompile(`^\d{7,15}$`)
	return phoneRegex.MatchString(phone)
}

// ValidateName checks if the name is valid
func ValidateName(name string) (bool, []string) {
	var errors []string

	// Trim whitespace
	name = strings.TrimSpace(name)

	// Check length
	if len(name) < 1 {
		errors = append(errors, "Name cannot be empty")
	}
	if len(name) > 50 {
		errors = append(errors, "Name must not exceed 50 characters")
	}

	// Check allowed characters (letters, spaces, hyphens, apostrophes)
	nameRegex := regexp.MustCompile(`^[a-zA-ZÀ-ÿ\s'-]+$`)
	if !nameRegex.MatchString(name) {
		errors = append(errors, "Name can only contain letters, spaces, hyphens, and apostrophes")
	}

	return len(errors) == 0, errors
}

// ValidateRequired checks if a string field is not empty
func ValidateRequired(value, fieldName string) (bool, string) {
	value = strings.TrimSpace(value)
	if value == "" {
		return false, fieldName + " is required"
	}
	return true, ""
}

// ValidateStringLength checks if a string length is within the specified range
func ValidateStringLength(value string, min, max int, fieldName string) (bool, []string) {
	var errors []string

	if len(value) < min {
		errors = append(errors, fieldName+" must be at least "+string(rune(min))+" characters long")
	}
	if len(value) > max {
		errors = append(errors, fieldName+" must not exceed "+string(rune(max))+" characters")
	}

	return len(errors) == 0, errors
}

// ValidateRole checks if the role is valid
func ValidateRole(role string) bool {
	validRoles := []string{"user", "admin", "moderator", "shipper", "shop_owner"}
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

// SanitizeString removes potentially harmful characters and trims whitespace
func SanitizeString(input string) string {
	// Trim whitespace
	input = strings.TrimSpace(input)

	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")

	// Remove other control characters if needed
	// input = strings.Map(func(r rune) rune {
	//     if unicode.IsControl(r) && r != '\t' && r != '\n' && r != '\r' {
	//         return -1
	//     }
	//     return r
	// }, input)

	return input
}
