package main

import (
	"log"
	"log-service/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {

	log.Println("logger service entering WriteLog() ...") //debugging

	//read json into a var
	var requestPayload JSONPayload
	_ = app.readJSON(w, r, &requestPayload)

	//insert data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, err)
		log.Println("WriteLog() inserted log in JSON format ...") //debugging
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	log.Println("WriteLog() preparing to write log ...") //debugging

	app.writeJSON(w, http.StatusAccepted, resp)

	log.Println("WriteLog() wrote logging status ...") //debugging

}
