package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {

	log.Println("entering Authenticate() ...") //debugging

	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)

	log.Println("read the request payload Authenticate() ...") //debugging

	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		log.Println("Bad request while authenticating (readJSON) ...") //debugging
		return
	}

	//validate the user against the database
	user, err := app.Models.User.GetByEmail(requestPayload.Email)

	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		log.Println("invalid credentials (GetByEmail) ...") //debugging
		return
	}

	log.Println("successfully validated user against pg database Authenticate() ...") //debugging

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		log.Println("invalid credentials ...") //debugging
		return
	}

	//log the authentication
	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		app.errorJSON(w, err)
		log.Println("failed to authenticate (logRequest) ...") //debugging
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)

}

func (app *Config) logRequest(name, data string) error {

	log.Println("Entering logRequest() ...") //debugging

	var entry struct {
		Name string `json: "name"`
		Data string `json: "data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceURL := "http://logger-service/log"

	log.Println("got log service URL (logRequest) ...", logServiceURL) //debugging

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("doing NewRequest() ...") //debugging
		return err
	}

	log.Println("got request from http.NewRequest (logRequest) ...", request) //debugging

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	log.Println("returning from logRequest() ...") //debugging

	return nil
}

// A basic  health checl function modelled closely on func "Authenticate"
func (app *Config) healthCheck(w http.ResponseWriter, r *http.Request) {

	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "I am Alive!",
		Data:    "you have reached the authentication service on the port 8081 -> 80",
	}

	app.writeJSON(w, http.StatusAccepted, payload)

}
