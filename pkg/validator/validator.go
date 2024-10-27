package validator

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Rule    string
	Message string
}

// CustomValidationFunc is a type for custom validation functions
type CustomValidationFunc func(interface{}) bool

// Validator represents the main validator struct
type Validator struct {
	errors           []ValidationError
	customValidators map[string]CustomValidationFunc
}

// New creates a new validator instance
func New() *Validator {
	return &Validator{
		errors:           make([]ValidationError, 0),
		customValidators: make(map[string]CustomValidationFunc),
	}
}

// RegisterCustomValidation registers a custom validation function
func (v *Validator) RegisterCustomValidation(name string, fn CustomValidationFunc) {
	v.customValidators[name] = fn
}

// GetErrors returns all validation errors
func (v *Validator) GetErrors() []ValidationError {
	return v.errors
}

// Validate performs validation on the given struct
func (v *Validator) Validate(s interface{}) error {
	v.errors = []ValidationError{} // Reset errors
	val := reflect.ValueOf(s)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return errors.New("validation only works on structs")
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		validateTag := fieldType.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		rules := strings.Split(validateTag, ",")
		for _, rule := range rules {
			v.validateField(fieldType.Name, field.Interface(), rule)
		}
	}

	if len(v.errors) > 0 {
		return errors.New("validation failed")
	}

	return nil
}

// addError adds a validation error
func (v *Validator) addError(field, rule, message string) {
	v.errors = append(v.errors, ValidationError{
		Field:   field,
		Rule:    rule,
		Message: message,
	})
}

// validateField validates a single field against a rule
func (v *Validator) validateField(fieldName string, value interface{}, rule string) {
	parts := strings.Split(rule, "=")
	ruleName := parts[0]
	ruleValue := ""
	if len(parts) > 1 {
		ruleValue = parts[1]
	}

	switch ruleName {
	// Basic validations
	case "required":
		if !v.required(value) {
			v.addError(fieldName, ruleName, "field is required")
		}
	case "notnil":
		if value == nil {
			v.addError(fieldName, ruleName, "field cannot be nil")
		}

	// String validations
	case "email":
		str, ok := value.(string)
		if !ok {
			v.addError(fieldName, ruleName, "field must be a string")
			return
		}
		if !v.email(str) {
			v.addError(fieldName, ruleName, "invalid email format")
		}
	case "url":
		str, ok := value.(string)
		if !ok {
			v.addError(fieldName, ruleName, "field must be a string")
			return
		}
		if !v.url(str) {
			v.addError(fieldName, ruleName, "invalid URL format")
		}
	case "alpha":
		str, ok := value.(string)
		if !ok {
			v.addError(fieldName, ruleName, "field must be a string")
			return
		}
		if !v.alpha(str) {
			v.addError(fieldName, ruleName, "field must contain only letters")
		}
	case "alphanum":
		str, ok := value.(string)
		if !ok {
			v.addError(fieldName, ruleName, "field must be a string")
			return
		}
		if !v.alphanum(str) {
			v.addError(fieldName, ruleName, "field must contain only letters and numbers")
		}
	case "numeric":
		str, ok := value.(string)
		if !ok {
			v.addError(fieldName, ruleName, "field must be a string")
			return
		}
		if !v.numeric(str) {
			v.addError(fieldName, ruleName, "field must contain only numbers")
		}
	case "lowercase":
		str, ok := value.(string)
		if !ok {
			v.addError(fieldName, ruleName, "field must be a string")
			return
		}
		if !v.lowercase(str) {
			v.addError(fieldName, ruleName, "field must be lowercase")
		}
	case "uppercase":
		str, ok := value.(string)
		if !ok {
			v.addError(fieldName, ruleName, "field must be a string")
			return
		}
		if !v.uppercase(str) {
			v.addError(fieldName, ruleName, "field must be uppercase")
		}

	// Length validations
	case "min":
		v.validateMin(fieldName, value, ruleValue)
	case "max":
		v.validateMax(fieldName, value, ruleValue)
	case "len":
		v.validateLen(fieldName, value, ruleValue)

	// Number range validations
	case "range":
		v.validateRange(fieldName, value, ruleValue)

	// Pattern validation
	case "pattern":
		str, ok := value.(string)
		if !ok {
			v.addError(fieldName, ruleName, "field must be a string")
			return
		}
		if !v.pattern(str, ruleValue) {
			v.addError(fieldName, ruleName, "field does not match pattern")
		}

	// Time validations
	case "datetime":
		str, ok := value.(string)
		if !ok {
			v.addError(fieldName, ruleName, "field must be a string")
			return
		}
		if !v.datetime(str, ruleValue) {
			v.addError(fieldName, ruleName, "invalid datetime format")
		}
	case "future":
		t, ok := value.(time.Time)
		if !ok {
			v.addError(fieldName, ruleName, "field must be a time.Time")
			return
		}
		if !v.future(t) {
			v.addError(fieldName, ruleName, "time must be in the future")
		}
	case "past":
		t, ok := value.(time.Time)
		if !ok {
			v.addError(fieldName, ruleName, "field must be a time.Time")
			return
		}
		if !v.past(t) {
			v.addError(fieldName, ruleName, "time must be in the past")
		}

	// Slice validations
	case "unique":
		if !v.unique(value) {
			v.addError(fieldName, ruleName, "slice must contain unique values")
		}

	// Custom validation
	default:
		if fn, ok := v.customValidators[ruleName]; ok {
			if !fn(value) {
				v.addError(fieldName, ruleName, fmt.Sprintf("failed custom validation: %s", ruleName))
			}
		}
	}
}

// Validation implementations
func (v *Validator) required(value interface{}) bool {
	if value == nil {
		return false
	}

	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.String:
		return strings.TrimSpace(val.String()) != ""
	case reflect.Slice, reflect.Map, reflect.Array:
		return val.Len() > 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return val.Uint() != 0
	case reflect.Float32, reflect.Float64:
		return val.Float() != 0
	case reflect.Bool:
		return val.Bool()
	case reflect.Ptr, reflect.Interface:
		return !val.IsNil()
	}
	return true
}

func (v *Validator) email(value string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, value)
	return match
}

func (v *Validator) url(value string) bool {
	_, err := url.ParseRequestURI(value)
	return err == nil
}

func (v *Validator) alpha(value string) bool {
	for _, r := range value {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func (v *Validator) alphanum(value string) bool {
	for _, r := range value {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}

func (v *Validator) numeric(value string) bool {
	for _, r := range value {
		if !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}

func (v *Validator) lowercase(value string) bool {
	return strings.ToLower(value) == value
}

func (v *Validator) uppercase(value string) bool {
	return strings.ToUpper(value) == value
}

func (v *Validator) validateMin(fieldName string, value interface{}, minStr string) {
	min, err := strconv.Atoi(minStr)
	if err != nil {
		v.addError(fieldName, "min", "invalid min value")
		return
	}

	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.String:
		if len(val.String()) < min {
			v.addError(fieldName, "min", fmt.Sprintf("length must be at least %d", min))
		}
	case reflect.Slice, reflect.Map, reflect.Array:
		if val.Len() < min {
			v.addError(fieldName, "min", fmt.Sprintf("length must be at least %d", min))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val.Int() < int64(min) {
			v.addError(fieldName, "min", fmt.Sprintf("value must be at least %d", min))
		}
	case reflect.Float32, reflect.Float64:
		if val.Float() < float64(min) {
			v.addError(fieldName, "min", fmt.Sprintf("value must be at least %d", min))
		}
	}
}

func (v *Validator) validateMax(fieldName string, value interface{}, maxStr string) {
	max, err := strconv.Atoi(maxStr)
	if err != nil {
		v.addError(fieldName, "max", "invalid max value")
		return
	}

	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.String:
		if len(val.String()) > max {
			v.addError(fieldName, "max", fmt.Sprintf("length must not exceed %d", max))
		}
	case reflect.Slice, reflect.Map, reflect.Array:
		if val.Len() > max {
			v.addError(fieldName, "max", fmt.Sprintf("length must not exceed %d", max))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val.Int() > int64(max) {
			v.addError(fieldName, "max", fmt.Sprintf("value must not exceed %d", max))
		}
	case reflect.Float32, reflect.Float64:
		if val.Float() > float64(max) {
			v.addError(fieldName, "max", fmt.Sprintf("value must not exceed %d", max))
		}
	}
}

func (v *Validator) validateLen(fieldName string, value interface{}, lenStr string) {
	length, err := strconv.Atoi(lenStr)
	if err != nil {
		v.addError(fieldName, "len", "invalid length value")
		return
	}

	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.String:
		if len(val.String()) != length {
			v.addError(fieldName, "len", fmt.Sprintf("length must be exactly %d", length))
		}
	case reflect.Slice, reflect.Map, reflect.Array:
		if val.Len() != length {
			v.addError(fieldName, "len", fmt.Sprintf("length must be exactly %d", length))
		}
	}
}

func (v *Validator) validateRange(fieldName string, value interface{}, rangeStr string) {
	parts := strings.Split(rangeStr, "-")
	if len(parts) != 2 {
		v.addError(fieldName, "range", "invalid range format")
		return
	}

	min, err1 := strconv.ParseFloat(parts[0], 64)
	max, err2 := strconv.ParseFloat(parts[1], 64)
	if err1 != nil || err2 != nil {
		v.addError(fieldName, "range", "invalid range values")
		return
	}

	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		num := float64(val.Int())
		if num < min || num > max {
			v.addError(fieldName, "range", fmt.Sprintf("value must be between %v and %v", min, max))
		}
	case reflect.Float32, reflect.Float64:
		num := val.Float()
		if num < min || num > max {
			v.addError(fieldName, "range", fmt.Sprintf("value must be between %v and %v", min, max))
		}
	}
}

func (v *Validator) pattern(value string, pattern string) bool {
	match, err := regexp.MatchString(pattern, value)
	return err == nil && match
}

func (v *Validator) datetime(value string, layout string) bool {
	if layout == "" {
		layout = time.RFC3339
	}
	_, err := time.Parse(layout, value)
	return err == nil
}

func (v *Validator) future(t time.Time) bool {
	return t.After(time.Now())
}

func (v *Validator) past(t time.Time) bool {
	return t.Before(time.Now())
}

func (v *Validator) unique(value interface{}) bool {
	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Slice {
		return false
	}

	seen := make(map[interface{}]bool)
	for i := 0; i < val.Len(); i++ {
		item := val.Index(i).Interface()
		if seen[item] {
			return false
		}
		seen[item] = true
	}
	return true
}
