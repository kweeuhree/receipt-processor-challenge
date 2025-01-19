package handlers

import "kweeuhree.receipt-processor-challenge/internal/models"

var SimpleReceipt = &models.Receipt{
	ID:           "123-qwe-456-rty-7890",
	Retailer:     "Target",
	PurchaseDate: "2022-01-01",
	PurchaseTime: "13:01",
	Items: []models.Item{
		{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
		{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
		{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
		{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
		{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
	},
	Total:  "35.35",
	Points: 28,
}

var ValidReceipt = &ReceiptInput{
	Retailer:     "Target",
	PurchaseDate: "2022-01-02",
	PurchaseTime: "13:13",
	Total:        "1.25",
	Items: []models.Item{
		{ShortDescription: "Pepsi - 12-oz", Price: "1.25"},
	},
}

var NoRetailerReceipt = &ReceiptInput{
	PurchaseDate: "2022-01-02",
	PurchaseTime: "13:13",
	Total:        "1.25",
	Items: []models.Item{
		{ShortDescription: "Pepsi - 12-oz", Price: "1.25"},
	},
}

var NoTotalReceipt = &ReceiptInput{
	Retailer:     "Target",
	PurchaseDate: "2022-01-01",
	PurchaseTime: "12:00",
	Items: []models.Item{
		{ShortDescription: "Pepsi - 12-oz", Price: "1.25"},
	},
}

var InvalidTotalReceipt = &ReceiptInput{
	Retailer:     "Target",
	PurchaseDate: "2022-01-01",
	PurchaseTime: "12:00",
	Total:        "hello-world",
	Items: []models.Item{
		{ShortDescription: "Pepsi - 12-oz", Price: "1.25"},
	},
}

var NoItemsReceipt = &ReceiptInput{
	Retailer:     "Target",
	PurchaseDate: "2022-01-01",
	PurchaseTime: "12:00",
	Total:        "1.25",
}
