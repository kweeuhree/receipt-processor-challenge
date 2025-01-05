package models

import (
	"testing"
)

var SimpleReceipt = &Receipt{
	ID:           "123-qwe-456-rty-7890",
	Retailer:     "Target",
	PurchaseDate: "2022-01-02",
	PurchaseTime: "13:13",
	Total:        "1.25",
	Items: []Item{
		{ShortDescription: "Pepsi - 12-oz", Price: "1.25"},
	},
}
var store *ReceiptStore

func TestMain(m *testing.M) {
	store = NewStore()
	m.Run()

}

func TestNewStore(t *testing.T) {
	if store.receipts == nil {
		t.Errorf("Expected receipts field to be a map, but got nil")
	}
}

func TestInsert(t *testing.T) {
	tests := []struct {
		name    string
		receipt *Receipt
	}{
		{"Insert empty receipt", &Receipt{}},
		{"Insert non-empty receipt", SimpleReceipt},
	}
	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			_ = store.Insert(*entry.receipt)

			if entry.receipt == nil && len(store.receipts) == 1 {
				t.Errorf("Expected receipts map to be empty, but got length %d", len(store.receipts))
			}

			if entry.receipt != nil && len(store.receipts) != 1 {
				t.Errorf("Expected a receipt, but got length %d", len(store.receipts))
			}

			_, exists := store.receipts[entry.receipt.ID]
			if !exists {
				t.Errorf("receipt with ID %v was not inserted", entry.receipt.ID)
			}

			t.Cleanup(func() {
				store = NewStore()
			})
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name    string
		receipt *Receipt
	}{
		{"No receipt to get", &Receipt{}},
		{"Valid receipt with an ID to get", SimpleReceipt},
	}
	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			_ = store.Insert(*entry.receipt)
			inserted, err := store.Get(entry.receipt.ID)
			if err != nil {
				t.Errorf("Could not get receipt with ID %s", entry.receipt.ID)
			}

			if inserted.ID != entry.receipt.ID {
				t.Errorf("Expected ID %s, received %s", entry.receipt.ID, inserted.ID)
			}

			t.Cleanup(func() {
				store = NewStore()
			})
		})
	}
}

func TestGetInvalidID(t *testing.T) {
	_, err := store.Get("invalid-id")
	if err == nil {
		t.Errorf("Expected an error, but got none")
	}

	expectedError := "no receipt found for that ID"
	if err.Error() != expectedError {
		t.Errorf("Expected an error message '%s', but got '%s'", expectedError, err.Error())
	}
}
