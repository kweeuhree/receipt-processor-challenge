package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter" // Third-party router for lightweight, efficient HTTP routing
	"github.com/justinas/alice"           // Middleware chaining library for clean, reusable middleware
)

// Initializes and configures the application's HTTP routes and middleware chain
func (app *application) routes() http.Handler {
	// Initialize the router
	router := httprouter.New()

	// Get receipt id
	router.Handler(http.MethodPost, "/receipts/process", http.HandlerFunc(app.handlers.ProcessReceipt))

	// Get receipt points
	router.Handler(http.MethodGet, "/receipts/:id/points", http.HandlerFunc(app.handlers.GetReceiptPoints))

	// Initialize the middleware chain using alice
	// Includes:
	// - recoverPanic: Middleware to recover from panics and prevent server crashes;
	// - logRequest: Middleware to log incoming HTTP requests.
	standard := alice.New(app.recoverPanic, app.logRequest)

	// Return the 'standard' middleware chain
	return standard.Then(router)
}
