package handlers

import (
	"kweeuhree.receipt-processor-challenge/internal/validator"
)

var v *validator.Validator

func (input *ReceiptInput) Validate() {
	input.CheckField(v.NotBlank(input.Retailer), "retailerName", "This field cannot be blank")
	input.CheckField(v.NotBlank(input.PurchaseDate), "purchaseDate", "This field cannot be blank")
	input.CheckField(v.ValidDate(input.PurchaseDate), "purchaseDate", "This field must be a valid date")
	input.CheckField(v.ValidTime(input.PurchaseTime), "purchaseTime", "This field must be valid time")
	input.CheckField(v.NotBlank(input.PurchaseTime), "purchaseTime", "This field cannot be blank")
	input.CheckField(v.NotBlank(input.Total), "total", "This field cannot be blank")
	input.CheckField(v.ValidNumber(input.Total), "total", "This field must be a valid number")
	input.CheckField(v.ItemsNotEmpty(input.Items), "items", "This field must have at least one object")
	for _, item := range input.Items {
		input.CheckField(v.ValidNumber(item.Price), "items", "Each item price must be a valid number")
	}
}
