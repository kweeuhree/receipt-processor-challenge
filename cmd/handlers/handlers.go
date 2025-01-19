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
	InfoLog      *log.Logger
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

func NewHandlers(errorLog *log.Logger, infoLog *log.Logger, receiptStore *models.ReceiptStore, utils *utils.Utils, helpers *helpers.Helpers) *Handlers {
	return &Handlers{
		ErrorLog:     errorLog,
		InfoLog:      infoLog,
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

	// Create and store new receipt
	newReceiptID, err := h.CreateAndStore(input)

	if err != nil {
		h.ErrorLog.Printf("Failed to store receipt: %v", err)
		h.Helpers.ServerError(w, err)
		return
	}

	// Construct the response
	response := IdResponse{
		ID: newReceiptID,
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

	// Get receipt by its id
	receipt, err := h.ReceiptStore.Get(receiptID)

	if err != nil {
		msg := map[string]string{"error": "No receipt found for that ID."}
		h.Helpers.EncodeJSON(w, http.StatusNotFound, msg)
		return
	}

	// Construct the response
	response := PointsResponse{
		Points: receipt.Points,
	}

	// Write the response struct to the response as JSON
	err = h.Helpers.EncodeJSON(w, http.StatusOK, response)
	if err != nil {
		h.Helpers.ServerError(w, err)
		return
	}
}

func (h *Handlers) CreateAndStore(input ReceiptInput) (string, error) {
	// Prepare new receipt for storage
	newReceipt, err := h.ReceiptFactory(input)
	if err != nil {
		return "", err
	}

	// Store the receipt in memory
	err = h.ReceiptStore.Insert(newReceipt)
	if err != nil {
		return "", err
	}

	return newReceipt.ID, nil
}

// Constructs a new receipt based on the input
func (h *Handlers) ReceiptFactory(input ReceiptInput) (models.Receipt, error) {
	receiptID := uuid.New().String()

	h.InfoLog.Printf("Calculating points for receipt with id: %s", receiptID)

	points, err := h.Utils.CalculatePoints(input.Retailer, input.PurchaseDate, input.PurchaseTime, input.Total, input.Items)
	if err != nil {
		return models.Receipt{}, err
	}

	h.InfoLog.Printf("Total Points: %d", points)

	newReceipt := models.Receipt{
		ID:           receiptID,
		Retailer:     input.Retailer,
		PurchaseDate: input.PurchaseDate,
		PurchaseTime: input.PurchaseTime,
		Total:        input.Total,
		Items:        input.Items,
		Points:       points,
	}

	return newReceipt, nil
}

func (h *Handlers) DeleteReceipt(w http.ResponseWriter, r *http.Request) {
	receiptID := h.Helpers.GetIdFromParams(r, "id")

	// Remove the receipt from receiptStore
	err := h.ReceiptStore.Delete(receiptID)

	if err != nil {
		h.ErrorLog.Printf("Failed to delete the receipt with ID %s. Error: %+v", receiptID, err)
		return
	}
	h.Helpers.EncodeJSON(w, http.StatusNoContent, "")
}
