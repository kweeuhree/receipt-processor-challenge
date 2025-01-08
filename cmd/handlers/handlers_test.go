package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"kweeuhree.receipt-processor-challenge/cmd/helpers"
	"kweeuhree.receipt-processor-challenge/cmd/utils"
	"kweeuhree.receipt-processor-challenge/internal/models"
)

type TestDependencies struct {
	receiptStore *models.ReceiptStore
	utils        *utils.Utils
	helpers      *helpers.Helpers
	handlers     *Handlers
}

func setupTestDependencies() *TestDependencies {
	receiptStore := models.NewStore()
	utils := &utils.Utils{}
	helpers := &helpers.Helpers{}
	handlers := NewHandlers(log.Default(), log.Default(), receiptStore, utils, helpers)

	return &TestDependencies{
		receiptStore: receiptStore,
		utils:        utils,
		helpers:      helpers,
		handlers:     handlers,
	}
}

func TestProcessReceipt(t *testing.T) {
	deps := setupTestDependencies()
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
			req := httptest.NewRequest(http.MethodPost, "/receipts/process", bytes.NewBuffer(body))
			// req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			// Test ProcessReceipt function
			deps.handlers.ProcessReceipt(resp, req)

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

		})
	}
}

func TestGetReceiptPoints(t *testing.T) {
	router := httprouter.New()
	// Register the route
	router.GET("/receipts/:id/points", mockHandler)

	tests := []struct {
		name   string
		ID     string
		url    string
		status int
	}{
		{"Valid receipt id", SimpleReceipt.ID, "/receipts/123-qwe-456-rty-7890/points", http.StatusOK},
		{"Invalid receipt id", "hello-world", "/receipts/hello-world/points", http.StatusNotFound},
		{"Invalid request url", SimpleReceipt.ID, "/hello/123-qwe-456-rty-7890/world", http.StatusNotFound},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			receiptStore.Insert(*SimpleReceipt)

			// Create request and response
			req := httptest.NewRequest(http.MethodGet, entry.url, nil)
			// Create request context that will enable id extraction
			params := httprouter.Params{
				httprouter.Param{Key: "id", Value: entry.ID},
			}
			ctx := context.WithValue(req.Context(), httprouter.ParamsKey, params)
			req = req.WithContext(ctx)
			// Create response
			resp := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(resp, req)

			// Test GetReceiptPoints function
			handlers.GetReceiptPoints(resp, req)

			// Check response status
			if resp.Code == http.StatusOK && entry.status == http.StatusOK {
				var response PointsResponse
				// Decode the request body
				if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				// Check if the points match
				if response.Points != SimpleReceipt.Points {
					t.Fatalf("Expected %d points, received %d", SimpleReceipt.Points, response.Points)
				}
			}

			t.Cleanup(func() {
				receiptStore = models.NewStore()
			})
		})
	}
}

func TestGetReceiptPoints(t *testing.T) {
	router := httprouter.New()
	// Register the route
	router.GET("/receipts/:id/points", mockHandler)

	tests := []struct {
		name   string
		ID     string
		url    string
		status int
	}{
		{"Valid receipt id", SimpleReceipt.ID, "/receipts/123-qwe-456-rty-7890/points", http.StatusOK},
		{"Invalid receipt id", "hello-world", "/receipts/hello-world/points", http.StatusNotFound},
		{"Invalid request url", SimpleReceipt.ID, "/hello/123-qwe-456-rty-7890/world", http.StatusNotFound},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			receiptStore.Insert(*SimpleReceipt)

			// Create request and response
			req := httptest.NewRequest(http.MethodGet, entry.url, nil)
			// Create request context that will enable id extraction
			params := httprouter.Params{
				httprouter.Param{Key: "id", Value: entry.ID},
			}
			ctx := context.WithValue(req.Context(), httprouter.ParamsKey, params)
			req = req.WithContext(ctx)
			// Create response
			resp := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(resp, req)

			// Test GetReceiptPoints function
			handlers.GetReceiptPoints(resp, req)

			// Check response status
			if resp.Code == http.StatusOK && entry.status == http.StatusOK {
				var response PointsResponse
				// Decode the request body
				if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				// Check if the points match
				if response.Points != SimpleReceipt.Points {
					t.Fatalf("Expected %d points, received %d", SimpleReceipt.Points, response.Points)
				}
			}

			t.Cleanup(func() {
				receiptStore = models.NewStore()
			})
		})
	}
}
