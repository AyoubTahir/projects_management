package validator

/*import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Validator represents the validator object
type Validator struct {
	Errors []error
}

// New creates a new validator instance
func NewValidator() *Validator {
	return &Validator{
		Errors: make([]error, 0),
	}
}

// Valid returns true if there are no errors
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error to the validator
func (v *Validator) AddError(field, message string) {
	v.Errors = append(v.Errors, &ValidationError{
		Field:   field,
		Message: message,
	})
}

// GetErrors returns all validation errors
func (v *Validator) GetErrors() []error {
	return v.Errors
}

// Required checks if a string value is not empty
func (v *Validator) Required(field, value string) *Validator {
	if strings.TrimSpace(value) == "" {
		v.AddError(field, "this field is required")
	}
	return v
}

// MinLength checks if a string value meets the minimum length
func (v *Validator) MinLength(field, value string, minLength int) *Validator {
	if len(value) < minLength {
		v.AddError(field, fmt.Sprintf("minimum length is %d characters", minLength))
	}
	return v
}

// MaxLength checks if a string value meets the maximum length
func (v *Validator) MaxLength(field, value string, maxLength int) *Validator {
	if len(value) > maxLength {
		v.AddError(field, fmt.Sprintf("maximum length is %d characters", maxLength))
	}
	return v
}

// Email validates email format
func (v *Validator) Email(field, email string) *Validator {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if matched, _ := regexp.MatchString(emailRegex, email); !matched && email != "" {
		v.AddError(field, "invalid email format")
	}
	return v
}

// Numeric checks if a string contains only numbers
func (v *Validator) Numeric(field, value string) *Validator {
	if value != "" {
		for _, char := range value {
			if !unicode.IsDigit(char) {
				v.AddError(field, "must contain only numbers")
				break
			}
		}
	}
	return v
}

// Range checks if a number is within a specific range
func (v *Validator) Range(field string, value, min, max int) *Validator {
	if value < min || value > max {
		v.AddError(field, fmt.Sprintf("must be between %d and %d", min, max))
	}
	return v
}

// Password validates password strength
func (v *Validator) Password(field, password string) *Validator {
	if password != "" {
		var (
			hasUpper   = false
			hasLower   = false
			hasNumber  = false
			hasSpecial = false
		)

		if len(password) < 8 {
			v.AddError(field, "password must be at least 8 characters long")
			return v
		}

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
			v.AddError(field, "password must contain at least one uppercase letter")
		}
		if !hasLower {
			v.AddError(field, "password must contain at least one lowercase letter")
		}
		if !hasNumber {
			v.AddError(field, "password must contain at least one number")
		}
		if !hasSpecial {
			v.AddError(field, "password must contain at least one special character")
		}
	}
	return v
}

// URL validates URL format
func (v *Validator) URL(field, url string) *Validator {
	urlRegex := `^(http|https):\/\/[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}([\/\w \.-]*)*\/?$`
	if matched, _ := regexp.MatchString(urlRegex, url); !matched && url != "" {
		v.AddError(field, "invalid URL format")
	}
	return v
}

// InList checks if a value is in a list of allowed values
func (v *Validator) InList(field, value string, list []string) *Validator {
	if value != "" {
		valid := false
		for _, item := range list {
			if value == item {
				valid = true
				break
			}
		}
		if !valid {
			v.AddError(field, fmt.Sprintf("must be one of: %s", strings.Join(list, ", ")))
		}
	}
	return v
}

// Custom allows adding custom validation rules
func (v *Validator) Custom(field string, validationFn func() error) *Validator {
	if err := validationFn(); err != nil {
		v.AddError(field, err.Error())
	}
	return v
}
*/
