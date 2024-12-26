package handlers

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"kweeuhree.receipt-processor-challenge/cmd/helpers"
	"kweeuhree.receipt-processor-challenge/cmd/utils"
	"kweeuhree.receipt-processor-challenge/internal/models"
	"kweeuhree.receipt-processor-challenge/internal/validator"
)

type Handlers struct {
	ErrorLog     *log.Logger
	ReceiptStore *models.ReceiptStore
	Utils        *utils.Utils
	Helpers      *helpers.Helpers
}

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

func NewHandlers(errorLog *log.Logger, receiptStore *models.ReceiptStore, utils *utils.Utils, helpers *helpers.Helpers) *Handlers {
	return &Handlers{
		ErrorLog:     errorLog,
		ReceiptStore: receiptStore,
		Utils:        utils,
		Helpers:      helpers,
	}
}

// Process the receipt and return its id
func (h *Handlers) ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON body into the input struct
	var input ReceiptInput
	err := h.Helpers.DecodeJSON(w, r, &input)
	if err != nil {
		h.ErrorLog.Printf("Exiting after decoding attempt: %s", err)
		return
	}

	// Validate input
	input.Validate()
	if !input.Valid() {
		h.Helpers.EncodeJSON(w, http.StatusBadRequest, input.FieldErrors)
		return
	}

	// Prepare new receipt for storage
	newReceipt := h.ReceiptFactory(input)

	// Store the receipt in memory
	err = h.ReceiptStore.Insert(newReceipt)
	if err != nil {
		h.Helpers.ServerError(w, err)
		return
	}

	// Construct the response
	response := IdResponse{
		ID: newReceipt.ID,
	}

	// Write the response struct as JSON
	err = h.Helpers.EncodeJSON(w, http.StatusOK, response)
	if err != nil {
		h.Helpers.ServerError(w, err)
		return
	}
}

// Process the receipt and return its id
func (h *Handlers) GetReceiptPoints(w http.ResponseWriter, r *http.Request) {
	// Get receipt id from params
	receiptID := h.Helpers.GetIdFromParams(r, "id")
	if receiptID == "" {
		h.Helpers.NotFound(w)
		return
	}

	// Calculate points
	points, err := h.Utils.CalculatePoints(receiptID)
	if err != nil {
		msg := map[string]string{"error": "No receipt found for that ID."}
		h.Helpers.EncodeJSON(w, http.StatusNotFound, msg)
		return
	}

	// Construct the response
	response := PointsResponse{
		Points: points,
	}

	// Write the response struct to the response as JSON
	err = h.Helpers.EncodeJSON(w, http.StatusOK, response)
	if err != nil {
		h.Helpers.ServerError(w, err)
		return
	}
}

func (h *Handlers) ReceiptFactory(input ReceiptInput) models.Receipt {
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
