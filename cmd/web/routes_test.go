package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
)

// Writes an OK status to the response
func mockHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}

func Test_routes(t *testing.T) {
	router := httprouter.New()

	// Register routes for testing
	router.POST("/receipts/process", mockHandler)
	router.GET("/receipts/:id/points", mockHandler)

	var registered = []struct {
		route          string
		method         string
		expectedStatus int
	}{
		{"/receipts/process", "POST", http.StatusOK},
		{"/receipts/123/points", "GET", http.StatusOK},
		{"/hello-world", "GET", http.StatusNotFound},
	}

	for _, route := range registered {
		// Check if the route exists
		if !routeExists(router, route.route, route.method, route.expectedStatus) {
			t.Errorf("route %s is not registered", route.route)
		}
	}
}

func routeExists(router *httprouter.Router, testRoute, testMethod string, expectedStatus int) bool {
	// Create a ResponseRecorder to capture the response
	recorder := httptest.NewRecorder()
	// Create a new HTTP request
	request := httptest.NewRequest(testMethod, testRoute, nil)
	// Serve the HTTP request
	router.ServeHTTP(recorder, request)

	return recorder.Code == expectedStatus
}
