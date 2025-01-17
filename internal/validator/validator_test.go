package validator

import (
	"testing"

	"kweeuhree.receipt-processor-challenge/internal/models"
	"kweeuhree.receipt-processor-challenge/testdata"
)

type TestValidator struct {
	validator *Validator
}

func setupTestDependencies() *TestValidator {
	v := &Validator{}
	return &TestValidator{validator: v}
}

func TestValid(t *testing.T) {
	tests := []struct {
		name      string
		validator Validator
		expected  bool
	}{
		{
			"Valid",
			Validator{},
			true,
		},
		{
			"Invalid with field errors",
			Validator{FieldErrors: map[string]string{"field": "error"}},
			false,
		},
		{
			"Invalid with non-field errors",
			Validator{NonFieldErrors: []string{"non-field error"}},
			false,
		},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			result := entry.validator.Valid()
			if result != entry.expected {
				t.Errorf("Expected %t, received %t", entry.expected, result)
			}
		})
	}
}

func TestAddNonFieldError(t *testing.T) {
	d := setupTestDependencies()
	tests := []struct {
		name string
		msg  string
	}{
		{
			"Empty message",
			" ",
		},
		{
			"Non-empty message",
			"non-field error",
		},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			// Reset NonFieldErrors before each test
			d.validator.NonFieldErrors = nil

			d.validator.AddNonFieldError(entry.msg)

			if d.validator.NonFieldErrors[0] != entry.msg {
				t.Errorf("Expected to receive %s, but did not", entry.msg)
			}
		})
	}
}

func TestAddFieldError(t *testing.T) {
	d := setupTestDependencies()
	tests := []struct {
		name string
		key  string
		msg  string
	}{
		{
			"Empty message with non-empty key",
			"key",
			" ",
		},
		{
			"Empty message with empty key",
			" ",
			" ",
		},
		{
			"Existing message and key",
			"key",
			"msg",
		},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			// Reset FieldErrors before each test
			d.validator.FieldErrors = nil

			d.validator.AddFieldError(entry.key, entry.msg)

			if d.validator.FieldErrors == nil {
				t.Errorf("Expected to receive key %s, but did not", entry.key)
			}

			if d.validator.FieldErrors[entry.key] != entry.msg {
				t.Errorf("Expected to receive message %s, but did not", entry.msg)
			}
		})
	}
}

func TestCheckField(t *testing.T) {
	d := setupTestDependencies()
	tests := []struct {
		name     string
		ok       bool
		key      string
		msg      string
		expected map[string]string
	}{
		{
			"Valid checkField",
			true,
			"field",
			"error",
			map[string]string{},
		},
		{
			"Invalid checkField",
			false,
			"field",
			"error",
			map[string]string{"field": "error"},
		},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			// Reset FieldErrors before each test
			d.validator.FieldErrors = nil

			d.validator.CheckField(entry.ok, entry.key, entry.msg)
			if len(d.validator.FieldErrors) != len(entry.expected) {
				t.Errorf("Expected %v errors, got %v", len(entry.expected), len(d.validator.FieldErrors))
			}
			for key, value := range entry.expected {
				if d.validator.FieldErrors[key] != value {
					t.Errorf("Expected value %v for key %v, got %v", value, key, d.validator.FieldErrors[key])
				}
			}
		})
	}
}

func TestNotBlank(t *testing.T) {
	d := setupTestDependencies()
	tests := []struct {
		name   string
		value  string
		result bool
	}{
		{
			"Non-empty value",
			"key",
			true,
		},
		{
			"Lots of spaces non-empty value",
			"   key    ",
			true,
		},
		{
			"Empty value",
			" ",
			false,
		},
		{
			"Lots of spaces empty value",
			"          ",
			false,
		},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			result := d.validator.NotBlank(entry.value)

			if result != entry.result {
				t.Errorf("Expected %t, but got %t", entry.result, result)
			}
		})
	}
}

func TestItemsNotEmpty(t *testing.T) {
	d := setupTestDependencies()
	tests := []struct {
		name   string
		items  []models.Item
		result bool
	}{
		{
			"Non-empty items",
			testdata.GatoradeReceiptItems,
			true,
		},
		{
			"Non-empty items",
			testdata.MountainDewReceiptItems,
			true,
		},
		{
			"Empty items",
			nil,
			false,
		},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			result := d.validator.ItemsNotEmpty(entry.items)

			if result != entry.result {
				t.Errorf("Expected %t, but got %t", entry.result, result)
			}
		})
	}
}

func TestValidDate(t *testing.T) {
	d := setupTestDependencies()
	tests := []struct {
		name   string
		date   string
		result bool
	}{
		{
			"Valid date",
			"2006-01-02",
			true,
		},
		{
			"Valid date",
			"1996-10-03",
			true,
		},
		{
			"Invalid date",
			"2016",
			false,
		},
		{
			"Invalid date",
			"2016-02-35",
			false,
		},
		{
			"Invalid date",
			"--02-",
			false,
		},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			result := d.validator.ValidDate(entry.date)

			if result != entry.result {
				t.Errorf("Expected %t, but got %t", entry.result, result)
			}
		})
	}
}

func TestValidTime(t *testing.T) {
	d := setupTestDependencies()
	tests := []struct {
		name       string
		timeString string
		result     bool
	}{
		{
			"Valid time",
			"15:14",
			true,
		},
		{
			"Valid time",
			"08:03",
			true,
		},
		{
			"Invalid time",
			"2016",
			false,
		},
		{
			"Invalid time",
			"88:03",
			false,
		},
		{
			"Invalid time",
			"02:",
			false,
		},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			result := d.validator.ValidTime(entry.timeString)

			if result != entry.result {
				t.Errorf("Expected %t, but got %t", entry.result, result)
			}
		})
	}
}
