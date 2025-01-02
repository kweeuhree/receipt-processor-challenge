package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"kweeuhree.receipt-processor-challenge/cmd/helpers"
)

// Declare application and logBuffer instance for all tests
var app *application
var logBuffer bytes.Buffer

func TestMain(m *testing.M) {
	// Initialize helpers struct
	helpers := &helpers.Helpers{
		ErrorLog: log.New(&logBuffer, "", log.LstdFlags),
	}
	// Initialize application struct
	app = &application{
		infoLog: log.New(&logBuffer, "", log.LstdFlags),
		helpers: helpers,
	}
	m.Run()
}

// Returns an http.HandlerFunc that either triggers a panic or writes an HTTP 200 status.
func testHandler(testPanic bool) http.HandlerFunc {
	fn := func(w http.ResponseWriter, req *http.Request) {
		if testPanic {
			panic("Test panic triggered")
		} else {
			w.WriteHeader(http.StatusOK)
		}

	}
	return http.HandlerFunc(fn)
}

// Ensures that logRequest middleware logs the correct request information
func Test_logRequest(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		url      string
		expected []string
	}{
		{
			name:     "GET",
			method:   http.MethodGet,
			url:      "/receipts/123/points",
			expected: []string{"HTTP/1", "GET", "/receipts/123/points"},
		},
		{
			name:     "POST",
			method:   http.MethodPost,
			url:      "/receipts/process",
			expected: []string{"HTTP/1", "POST", "/receipts/process"},
		},
		{
			name:     "Invalid URL",
			method:   http.MethodGet,
			url:      "/hello-world",
			expected: []string{"HTTP/1", "GET", "/hello-world"},
		},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			// Create middleware with the handler that does not trigger panic
			middleware := app.logRequest(testHandler(false))
			// Create a new HTTP request
			req, err := http.NewRequest(entry.method, entry.url, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Create a ResponseRecorder to capture the response
			resp := httptest.NewRecorder()

			// Serve the HTTP request
			middleware.ServeHTTP(resp, req)

			// Get log output
			logOutput := logBuffer.String()
			if logOutput == "" {
				t.Error("No information logged to the logger")
			}
			// Check the contents of the log output
			for _, part := range entry.expected {
				if !strings.Contains(logOutput, part) {
					t.Errorf("Expected log output to contain '%s', but it didn't. Log: %s", part, logOutput)
				}
			}
		})
	}
}

// Ensures that recoverPanic middleware handles panics and sets the correct response
func Test_recoverPanic(t *testing.T) {
	tests := []struct {
		name             string
		panicOccurred    bool
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:             "No panic",
			panicOccurred:    false,
			expectedStatus:   http.StatusOK,
			expectedResponse: "OK",
		},
		{
			name:             "Panic occurred",
			panicOccurred:    true,
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "Internal Server Error",
		},
	}

	for _, entry := range tests {
		t.Run(entry.name, func(t *testing.T) {
			// Create middleware with the handler that does not trigger panic
			middleware := app.recoverPanic(testHandler(entry.panicOccurred))
			// Create a new HTTP request
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Create a ResponseRecorder to capture the response
			resp := httptest.NewRecorder()

			// Serve the HTTP request
			middleware.ServeHTTP(resp, req)
			// Check if the received status code matches expected status
			if resp.Result().StatusCode != entry.expectedStatus {
				t.Errorf("Expected %d, got '%d'", entry.expectedStatus, resp.Result().StatusCode)
			}
			// Check if the "Connection" header is set to "close"
			if entry.panicOccurred {
				if connectionHeader := resp.Header().Get("Connection"); connectionHeader != "close" {
					t.Errorf("Expected 'Connection' header to be 'close', got '%s'", connectionHeader)
				}
			}
		})
	}
}
