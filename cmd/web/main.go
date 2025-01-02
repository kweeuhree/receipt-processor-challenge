package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"kweeuhree.receipt-processor-challenge/cmd/handlers"
	"kweeuhree.receipt-processor-challenge/cmd/helpers"
	"kweeuhree.receipt-processor-challenge/cmd/utils"
	"kweeuhree.receipt-processor-challenge/internal/models"
)

// Application-wide dependencies
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	// In-memory receipts storage
	receiptStore *models.ReceiptStore
	items        *[]models.Item
	handlers     *handlers.Handlers
	helpers      *helpers.Helpers
	utils        *utils.Utils
}

// Main point of entry
func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// Error and info logs
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	receiptStore := models.NewStore()
	utils := utils.NewUtils(errorLog, infoLog, receiptStore)
	helpers := helpers.NewHelpers(errorLog)
	handlers := handlers.NewHandlers(errorLog, receiptStore, utils, helpers)

	// Initialize the application with its dependencies
	app := &application{
		errorLog:     errorLog,
		infoLog:      infoLog,
		receiptStore: receiptStore,
		items:        &[]models.Item{},
		handlers:     handlers,
		utils:        utils,
		helpers:      helpers,
	}

	// HTTP server config
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Listen and serve
	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
