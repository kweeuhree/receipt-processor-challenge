package utils

import (
	"testing"

	"kweeuhree.receipt-processor-challenge/internal/models"
)

// Mock data for testing
var retailer = "Walgreens"
var purchaseDate = "2022-01-02"
var purchaseTime = "08:13"
var total = "2.65"
var items = []models.Item{
	{ShortDescription: "Pepsi - 12-oz", Price: "1.25"},
	{ShortDescription: "Dasani", Price: "1.40"},
}

// Sequential version benchmark
func BenchmarkCalculatePoints_Sequential(b *testing.B) {
	utils := Utils{}
	for i := 0; i < b.N; i++ {
		utils.CalculatePoints(retailer, purchaseDate, purchaseTime, total, items)
	}
}

// Concurrent version benchmark
func BenchmarkCalculatePoints_Concurrent(b *testing.B) {
	utils := Utils{}
	for i := 0; i < b.N; i++ {
		utils.ConcurrentCalculatePoints(retailer, purchaseDate, purchaseTime, total, items)
	}
}
