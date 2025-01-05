package handlers

import "kweeuhree.receipt-processor-challenge/internal/models"

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
