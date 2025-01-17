package models

import (
	"fmt"
)

type Receipt struct {
	ID           string
	Retailer     string
	PurchaseDate string
	PurchaseTime string
	Total        string
	Items        []Item
	Points       int
}

type ReceiptStore struct {
	receipts map[string]Receipt
}

func NewStore() *ReceiptStore {
	return &ReceiptStore{
		receipts: make(map[string]Receipt),
	}
}

func (s *ReceiptStore) Insert(receipt Receipt) error {
	s.receipts[receipt.ID] = receipt

	return nil
}

func (s *ReceiptStore) Get(id string) (Receipt, error) {
	receipt, exists := s.receipts[id]
	if !exists {
		return Receipt{}, fmt.Errorf("no receipt found for that ID")
	}

	return receipt, nil
}

func (s *ReceiptStore) Delete(id string) error {
	_, exists := s.receipts[id]
	if !exists {
		return fmt.Errorf("no receipt found for that ID")
	}

	delete(s.receipts, id)

	return nil
}
