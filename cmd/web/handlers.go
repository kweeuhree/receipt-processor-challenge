package main

import (
	"net/http"

	"github.com/google/uuid"
	"kweeuhree.receipt-processor-challenge/internal/models"
	"kweeuhree.receipt-processor-challenge/internal/validator"
)

type ReceiptInput struct {
	Retailer     string        `json:"retailer"`
	PurchaseDate string        `json:"purchaseDate"`
	PurchaseTime string        `json:"purchaseTime"`
	Total        string        `json:"total"`
	Items        []models.Item `json:"items"`
	validator.Validator
}

type IdResponse struct {
	ID string `json:"id"`
}

type PointsResponse struct {
	Points int `json:"points"`
}

// Process the receipt and return its id
func (app *application) ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON body into the input struct
	var input ReceiptInput
	err := decodeJSON(w, r, &input)
	if err != nil {
		app.errorLog.Printf("Exiting after decoding attempt: %s", err)
		return
	}

	// Validate input
	input.Validate()
	if !input.Valid() {
		encodeJSON(w, http.StatusBadRequest, input.FieldErrors)
		return
	}

	// Prepare new receipt for storage
	newReceipt := app.ReceiptFactory(input)

	// Store the receipt in memory
	err = app.receiptStore.Insert(newReceipt)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Construct the response
	response := IdResponse{
		ID: newReceipt.ID,
	}

	// Write the response struct as JSON
	err = encodeJSON(w, http.StatusOK, response)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

// Process the receipt and return its id
func (app *application) GetReceiptPoints(w http.ResponseWriter, r *http.Request) {
	// Get receipt id from params
	receiptID := app.GetIdFromParams(r, "id")
	if receiptID == "" {
		app.notFound(w)
		return
	}

	// Calculate points
	points, err := app.CalculatePoints(receiptID)
	if err != nil {
		msg := map[string]string{"error": "No receipt found for that ID."}
		encodeJSON(w, http.StatusNotFound, msg)
		return
	}

	// Construct the response
	response := PointsResponse{
		Points: points,
	}

	// Write the response struct to the response as JSON
	err = encodeJSON(w, http.StatusOK, response)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) ReceiptFactory(input ReceiptInput) models.Receipt {
	receiptID := uuid.New().String()

	newReceipt := models.Receipt{
		ID:           receiptID,
		Retailer:     input.Retailer,
		PurchaseDate: input.PurchaseDate,
		PurchaseTime: input.PurchaseTime,
		Total:        input.Total,
		Items:        input.Items,
	}

	return newReceipt
}
