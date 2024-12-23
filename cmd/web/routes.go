package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter" // router
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Initialize the router
	router := httprouter.New()

	// Get receipt id
	router.Handler(http.MethodPost, "/receipts/process", http.HandlerFunc(app.ProcessReceipt))

	// Get receipt points
	router.Handler(http.MethodGet, "/receipts/:id/points", http.HandlerFunc(app.GetReceiptPoints))

	// Chain middleware
	standard := alice.New(app.recoverPanic, app.logRequest)

	// Return the 'standard' middleware chain
	return standard.Then(router)
}
