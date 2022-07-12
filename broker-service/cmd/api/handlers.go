package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := JSONResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	payload := RequestPayload{}
	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.errorJSON(w, err)
	}

	switch payload.Action {
	case "auth":
		app.authenticate(w, payload.Auth)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}

}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create some json we'll send to the auth service
	data, _ := json.MarshalIndent(a, "", "\t")
	// call the service
	req, err := http.NewRequest("POST", "http://auth-service/auth", bytes.NewBuffer(data))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer res.Body.Close()

	// make sure we get back the correct status code
	if res.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if res.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	// create variable to read res.Body into
	dataFromService := &JSONResponse{}

	err = json.NewDecoder(res.Body).Decode(dataFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if dataFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	payload := JSONResponse{
		Error:   false,
		Message: "Authenticated!",
		Data:    dataFromService.Data,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
