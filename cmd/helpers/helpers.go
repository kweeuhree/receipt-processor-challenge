package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/julienschmidt/httprouter"
)

type Helpers struct {
	ErrorLog *log.Logger
}

func NewHelpers(errorLog *log.Logger) *Helpers {
	return &Helpers{
		ErrorLog: errorLog,
	}
}

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (h *Helpers) ServerError(w http.ResponseWriter, err error) {
	// Use the debug.Stack() function to get a stack trace for the current goroutine and append it to the log message
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// Report the file name and line number one step back in the stack trace
	// to have a clearer idea of where the error actually originated from
	// set frame depth to 2
	h.ErrorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description
// to the user
func (h *Helpers) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Not found helper
func (h *Helpers) NotFound(w http.ResponseWriter) {
	h.ClientError(w, http.StatusNotFound)
}

// Decode the JSON body of a request into the destination struct
func (h *Helpers) DecodeJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
		return err
	}
	return nil
}

// Encodes provided data into a JSON response
func (h *Helpers) EncodeJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// Get parameter id from the request
func (h *Helpers) GetIdFromParams(r *http.Request, paramsId string) string {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName(paramsId)
	return id
}
