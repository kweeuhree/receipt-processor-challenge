package main

import (
	"testing"

	"kweeuhree.receipt-processor-challenge/internal/models"
)

// Declare and initialize application instance for all tests
var app *application

func TestMain(m *testing.M) {
	app = &application{}
	m.Run()
}

func Test_getRetailerNamePoints(t *testing.T) {
	retailerNameTests := []struct {
		name     string
		expected int
	}{
		{"Target", 6},
		{"", 0},
		{"M&M Corner Market", 14},
	}

	for _, entry := range retailerNameTests {
		result := app.getRetailerNamePoints(entry.name)

		if result != entry.expected {
			t.Errorf("retailerNameTests. %s: expected %d, received %d", entry.name, entry.expected, result)
		}
	}
}

func Test_isAlphanumeric(t *testing.T) {
	isAlphanumericTests := []struct {
		char     string
		expected bool
	}{
		{"1", true},
		{"a", true},
		{"!", false},
		{"&", false},
		{" ", false},
	}

	for _, entry := range isAlphanumericTests {
		result := app.isAlphanumeric(entry.char)

		if result != entry.expected {
			t.Errorf("isAlphanumericTests. %s: expected %t, received %t", entry.char, entry.expected, result)
		}
	}
}

func Test_getRoundTotalPoints(t *testing.T) {
	roundTotalTests := []struct {
		num      float64
		expected int
	}{
		{0, 50},
		{1.00, 50},
		{15.00, 50},
		{99.99, 0},
		{1.10, 0},
	}

	for _, entry := range roundTotalTests {
		result := app.getRoundTotalPoints(entry.num)

		if result != entry.expected {
			t.Errorf("roundTotalTests. %f: expected %d, received %d", entry.num, entry.expected, result)
		}
	}
}

func Test_getQuartersPoints(t *testing.T) {
	quartersTests := []struct {
		num      float64
		expected int
	}{
		{0.25, 25},
		{1.00, 25},
		{15.00, 25},
		{99.99, 0},
		{1.10, 0},
	}

	for _, entry := range quartersTests {
		result := app.getQuartersPoints(entry.num)

		if result != entry.expected {
			t.Errorf("quartersTests. %f: expected %d, received %d", entry.num, entry.expected, result)
		}
	}
}

func Test_getEveryTwoItemsPoints(t *testing.T) {
	everyTwoItemsTests := []struct {
		items    []models.Item
		expected int
	}{
		{
			items: []models.Item{
				{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
				{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
				{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
				{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
				{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
			},
			expected: 10,
		},
		{
			items: []models.Item{
				{ShortDescription: "Gatorade", Price: "2.25"},
				{ShortDescription: "Gatorade", Price: "2.25"},
				{ShortDescription: "Gatorade", Price: "2.25"},
				{ShortDescription: "Gatorade", Price: "2.25"},
			},
			expected: 10,
		},
	}

	for index, entry := range everyTwoItemsTests {
		result := app.getEveryTwoItemsPoints(entry.items)

		if result != entry.expected {
			t.Errorf("everyTwoItemsTests[%d]: expected %d, received %d", index, entry.expected, result)
		}
	}
}

func Test_getOddDayPoints(t *testing.T) {
	oddDayTests := []struct {
		date     string
		expected int
	}{
		{"2022-01-01", 6},
		{"2022-02-15", 6},
		{"2022-03-21", 6},
		{"2022-04-10", 0},
		{"2022-05-26", 0},
	}

	for _, entry := range oddDayTests {
		result := app.getOddDayPoints(entry.date)

		if result != entry.expected {
			t.Errorf("oddDayTests. %s: expected %d, received %d", entry.date, entry.expected, result)
		}
	}
}

func Test_getPurchaseTimePoints(t *testing.T) {
	purchaseTimeTests := []struct {
		date     string
		expected int
	}{
		{"14:01", 10},
		{"14:23", 10},
		{"15:59", 10},
		{"16:00", 0},
		{"18:00", 0},
	}

	for _, entry := range purchaseTimeTests {
		result := app.getPurchaseTimePoints(entry.date)

		if result != entry.expected {
			t.Errorf("purchaseTimeTests. %s: expected %d, received %d", entry.date, entry.expected, result)
		}
	}
}
