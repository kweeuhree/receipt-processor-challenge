package models

import (
	"fmt"
)

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type Receipt struct {
	ID           *string `json:"id"`
	Retailer     string  `json:"retailer"`
	PurchaseDate string  `json:"purchaseDate"`
	PurchaseTime string  `json:"purchaseTime"`
	Total        string  `json:"total"`
	Items        []Item  `json:"items"`
}

type ReceiptStore struct {
	receipts map[string]Receipt
}

func NewStore() *ReceiptStore {
	return &ReceiptStore{
		receipts: make(map[string]Receipt),
	}
}

func (s *ReceiptStore) Insert(id string, receipt Receipt) error {
	receipt.ID = &id
	s.receipts[id] = receipt

	return nil
}

func (s *ReceiptStore) Get(id string) (Receipt, error) {
	receipt, exists := s.receipts[id]
	if !exists {
		return Receipt{}, fmt.Errorf("receipt not found")
	}

	return receipt, nil
}
