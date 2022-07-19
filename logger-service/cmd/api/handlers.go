package main

import (
	"net/http"

	"github.com/mrpineapples/logger/data"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// read json into var
	reqPayload := JSONPayload{}
	_ = app.readJSON(w, r, &reqPayload)

	// insert data
	event := data.LogEntry{
		Name: reqPayload.Name,
		Data: reqPayload.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	res := JSONResponse{
		Error:   false,
		Message: "logged",
	}

	app.writeJSON(w, http.StatusAccepted, res)
}
