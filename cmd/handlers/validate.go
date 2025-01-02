package handlers

import (
	"kweeuhree.receipt-processor-challenge/internal/validator"
)

func (input *ReceiptInput) Validate() {
	input.CheckField(validator.NotBlank(input.Retailer), "retailerName", "This field cannot be blank")
	input.CheckField(validator.NotBlank(input.PurchaseDate), "purchaseDate", "This field cannot be blank")
	input.CheckField(validator.ValidDate(input.PurchaseDate), "purchaseDate", "This field must be a valid date")
	input.CheckField(validator.NotBlank(input.PurchaseTime), "purchaseTime", "This field cannot be blank")
	input.CheckField(validator.NotBlank(input.Total), "total", "This field cannot be blank")
	input.CheckField(validator.ItemsNotEmpty(input.Items), "items", "This field must have at least one object")
}
