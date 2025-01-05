package validator

import (
	"strconv"
	"strings"
	"time"

	"kweeuhree.receipt-processor-challenge/internal/models"
)

// Validator struct contains a map of validation errors
type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

// Returns true if the FieldErrors and nonFieldErrors map doesn't contain any entries
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

// Adds an error message to the FieldErrors map
func (v *Validator) AddFieldError(key, message string) {
	// Initialize the map, if it isn't already initialized
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// Adds an error message to the FieldErrors map if a
// validation check is not 'ok'
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// Returns true if a value is not an empty string
func (v *Validator) NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// Returns true if a value is not an empty array
func (v *Validator) ItemsNotEmpty(items []models.Item) bool {
	return len(items) > 0
}

// Returns true if a value is a valid date
func (v *Validator) ValidDate(date string) bool {
	layout := "2006-01-02"
	_, err := time.Parse(layout, date)
	return err == nil
}

// Returns true if a value is valid time
func (v *Validator) ValidTime(timeString string) bool {
	layout := "15:04"
	_, err := time.Parse(layout, timeString)
	return err == nil
}

// Returns true if a total is valid number
func (v *Validator) ValidNumber(total string) bool {
	_, err := strconv.ParseFloat(total, 64)
	return err == nil
}
