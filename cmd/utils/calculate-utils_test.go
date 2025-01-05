package utils

import (
	"testing"

	"kweeuhree.receipt-processor-challenge/internal/models"
	"kweeuhree.receipt-processor-challenge/testdata"
)

// Declare and initialize application instance for all tests
var utils *Utils

func TestMain(m *testing.M) {
	utils = &Utils{}
	m.Run()
}

func Test_getRetailerNamePoints(t *testing.T) {
	tests := []struct {
		name     string
		expected int
	}{
		{"Target", 6},
		{"", 0},
		{"M&M Corner Market", 14},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			result := utils.getRetailerNamePoints(entry.name)

			if result != entry.expected {
				t.Errorf("Expected %d, received %d", entry.expected, result)
			}
		})
	}
}

func Test_isAlphanumeric(t *testing.T) {
	tests := []struct {
		name     string
		char     string
		expected bool
	}{
		{"Alphanumeric", "1", true},
		{"Alphanumeric", "a", true},
		{"Non-alphanumeric", "!", false},
		{"Non-alphanumeric", "&", false},
		{"Non-alphanumeric", " ", false},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			result := utils.isAlphanumeric(entry.char)

			if result != entry.expected {
				t.Errorf("For char %s: expected %t, received %t", entry.char, entry.expected, result)
			}
		})
	}
}

func Test_getRoundTotalPoints(t *testing.T) {
	tests := []struct {
		name     string
		num      float64
		expected int
	}{
		{"Points should be added", 0, 50},
		{"Points should be added", 1.00, 50},
		{"Points should be added", 15.00, 50},
		{"Points should not be added", 99.99, 0},
		{"Points should not be added", 1.10, 0},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			result := utils.getRoundTotalPoints(entry.num)

			if result != entry.expected {
				t.Errorf("For num %f: expected %d, received %d", entry.num, entry.expected, result)
			}
		})
	}
}

func Test_getQuartersPoints(t *testing.T) {
	tests := []struct {
		name     string
		num      float64
		expected int
	}{
		{"Points should be added", 0.25, 25},
		{"Points should be added", 1.00, 25},
		{"Points should be added", 15.00, 25},
		{"Points should not be added", 99.99, 0},
		{"Points should not be added", 1.10, 0},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			result := utils.getQuartersPoints(entry.num)

			if result != entry.expected {
				t.Errorf("For num %f: expected %d, received %d", entry.num, entry.expected, result)
			}
		})
	}
}

func Test_getEveryTwoItemsPoints(t *testing.T) {
	tests := []struct {
		name     string
		items    []models.Item
		expected int
	}{
		{"Mountain Dew receipt", testdata.MountainDewReceiptItems, 10},
		{"Gatorade receipt", testdata.GatoradeReceiptItems, 10},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			result := utils.getEveryTwoItemsPoints(entry.items)

			if result != entry.expected {
				t.Errorf("Expected %d, received %d", entry.expected, result)
			}
		})
	}
}

func Test_getItemDescriptionPoints(t *testing.T) {
	tests := []struct {
		name     string
		items    []models.Item
		expected int
	}{
		{"empty", nil, 0},
		{"Gatorade items", testdata.GatoradeReceiptItems, 0},
		{"Mountain Dew items", testdata.MountainDewReceiptItems, 6},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			result := utils.getItemDescriptionPoints(entry.items)

			if result != entry.expected {
				t.Errorf("Expected %d, received %d", entry.expected, result)
			}
		})
	}
}

func Test_getOddDayPoints(t *testing.T) {
	tests := []struct {
		name     string
		date     string
		expected int
	}{
		{"Odd day", "2022-01-01", 6},
		{"Odd day", "2022-02-15", 6},
		{"Odd day", "2022-03-21", 6},
		{"Even day", "2022-04-10", 0},
		{"Even day", "2022-05-26", 0},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			result := utils.getOddDayPoints(entry.date)

			if result != entry.expected {
				t.Errorf("For date %s: expected %d, received %d", entry.date, entry.expected, result)
			}
		})
	}
}

func Test_getPurchaseTimePoints(t *testing.T) {
	tests := []struct {
		date     string
		expected int
	}{
		{"14:01", 10},
		{"14:23", 10},
		{"15:59", 10},
		{"16:00", 0},
		{"18:00", 0},
	}

	for _, entry := range tests {
		t.Run(entry.date, func(t *testing.T) {
			result := utils.getPurchaseTimePoints(entry.date)

			if result != entry.expected {
				t.Errorf("Expected %d, received %d", entry.expected, result)
			}
		})
	}
}
