package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"kweeuhree.receipt-processor-challenge/cmd/helpers"
	"kweeuhree.receipt-processor-challenge/cmd/utils"
	"kweeuhree.receipt-processor-challenge/internal/models"
)

var handlers *Handlers
var receiptStore *models.ReceiptStore

func TestMain(m *testing.M) {
	receiptStore = models.NewStore()
	utils := &utils.Utils{}
	helpers := &helpers.Helpers{}
	handlers = NewHandlers(log.Default(), receiptStore, utils, helpers)
	m.Run()
}

func TestProcessReceipt(t *testing.T) {
	tests := []struct {
		name           string
		input          ReceiptInput
		expectedStatus int
		expectedID     string
		expectedField  string
		expectedError  string
	}{
		{"Valid receipt", *ValidReceipt, http.StatusOK, "123-qwe-456-rty-7890", "", ""},
		{"No retailer", *NoRetailerReceipt, http.StatusBadRequest, "", "retailerName", "This field cannot be blank"},
		{"No items", *NoItemsReceipt, http.StatusBadRequest, "", "items", "This field must have at least one object"},
		{"No total", *NoTotalReceipt, http.StatusBadRequest, "", "total", "This field cannot be blank"},
		{"Invalid total", *InvalidTotalReceipt, http.StatusBadRequest, "", "total", "This field must be a valid number"},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			// Prepare input for the request construction
			body, err := json.Marshal(entry.input)
			if err != nil {
				t.Fatalf("Failed to marshal input: %v", err)
			}
			// Create request and response
			req := httptest.NewRequest("POST", "/receipts/process", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			// Test ProcessReceipt function
			handlers.ProcessReceipt(resp, req)

			// Check the response status
			if status := resp.Code; status != entry.expectedStatus {
				t.Errorf("Expected status %d, got %d", entry.expectedStatus, status)
			}

			// Check if response status matches the expected one
			if entry.expectedStatus == http.StatusOK {
				// If so, check if the received id matches the expected one
				var response IdResponse
				if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if response.ID == "" {
					t.Errorf("Expected id %s, but did not receive any", entry.expectedID)
				}
			} else {
				// If response status does not match, check errors
				var response map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode error response: %v", err)
				}

				// Check both field and field error
				for field, fieldError := range response {
					if field != entry.expectedField {
						t.Errorf("Expected %s, received field %s", entry.expectedField, field)
					}
					if fieldError != entry.expectedError {
						t.Errorf("Expected %s, received field %s", entry.expectedError, fieldError)
					}
				}
			}

			t.Cleanup(func() {
				receiptStore = models.NewStore()
			})
		})
	}
}
