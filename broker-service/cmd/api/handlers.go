package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mrpineapples/broker/event"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	FromAddress string `json:"from"`
	To          string `json:"to"`
	Subject     string `json:"subject"`
	Message     string `json:"message"`
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
	case "log":
		app.logRabbitEvent(w, payload.Log)
	case "mail":
		app.sendMail(w, payload.Mail)
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

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	data, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"
	req, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(data))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		app.errorJSON(w, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		app.errorJSON(w, err)
	}

	payload := JSONResponse{
		Error:   false,
		Message: "logged",
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	data, _ := json.MarshalIndent(msg, "", "\t")

	mailServiceURL := "http://mail-service/send"
	req, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(data))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("Error calling mail service"))
		return
	}

	var payload JSONResponse
	payload.Error = false
	payload.Message = "Message sent to " + msg.To

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logRabbitEvent(w http.ResponseWriter, l LogPayload) {
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := JSONResponse{
		Message: "logged via RabbitMQ",
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	data, _ := json.MarshalIndent(payload, "", "\t")
	err = emitter.Push(string(data), "log.INFO")
	if err != nil {
		return err
	}

	return nil
}
