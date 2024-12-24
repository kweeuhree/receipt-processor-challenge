package main

import (
	"net/http"

	"github.com/google/uuid"
	"kweeuhree.receipt-processor-challenge/internal/models"
)

type IdResponse struct {
	ID string `json:"id"`
}

type PointsResponse struct {
	Points int `json:"points"`
}

// Process the receipt and return its id
func (app *application) ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON body into the input struct
	var input models.Receipt
	err := decodeJSON(w, r, &input)
	if err != nil {
		app.errorLog.Printf("Exiting after decoding attempt: %s", err)
		return
	}

	// Generate receipt id
	receiptID := uuid.New().String()

	// Store the receipt in memory
	err = app.receiptStore.Insert(receiptID, input)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Construct the response
	response := IdResponse{
		ID: receiptID,
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
		app.serverError(w, err)
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
