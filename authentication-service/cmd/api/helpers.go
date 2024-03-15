package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Read JSON
func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 //1 MB

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)

	if err != nil {
		log.Println("failed to read JSON (readJSON) ...") //debugging
		return err
	}

	err = dec.Decode(&struct{}{})

	if err != io.EOF {
		log.Println("failed to decode JSON (Decode) ...") //debugging
		return errors.New("body must have only a single JSON value")
	}

	log.Println("returning from readJSON() ...") //debugging

	return nil
}

// Write JSON
func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {

	out, err := json.Marshal(data)

	if err != nil {
		log.Println("ffailed to marshal JSON (writeJSON)") //debugging
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		log.Println("failed to write headers (writeJSON)") //debugging
		return err
	}

	log.Println("returning from writeJSON() ...") //debugging
	return nil
}

// Generate a JSON response
func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {

	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	log.Println("returning from errorJSON() ...") //debugging
	return app.writeJSON(w, statusCode, payload)
}
