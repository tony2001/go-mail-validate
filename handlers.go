package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"net/http"
	//"go-mail-validate/config"

	"github.com/tony2001/go-mail-validate/log"
	"github.com/tony2001/go-mail-validate/validate"
)

const RequestBodySizeLimit = 1024

type RequestEmail struct {
	Email string
}

func sendJsonResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Errorf("failed to send response: %s", err)
	}
}

//default handler
func handleDefault(w http.ResponseWriter, r *http.Request) {
	escapedURL := html.EscapeString(r.URL.Path)

	log.Infof("request to %q", escapedURL)

	errorStr := fmt.Sprintf("Non-existent endpoint %q", escapedURL)
	http.Error(w, errorStr, http.StatusNotFound)
}

func handleEmailValidate(w http.ResponseWriter, r *http.Request) {
	escapedURL := html.EscapeString(r.URL.Path)

	log.Infof("request to %q", escapedURL)

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, RequestBodySizeLimit))
	if err != nil {
		log.Warningf("failed to read request body: %s", err)
		response := NewResponseError("failed to read request data")
		sendJsonResponse(w, response)
		return
	}

	var reqEmail RequestEmail
	err = json.Unmarshal(body, &reqEmail)
	if err != nil {
		response := NewResponseError("invalid JSON request data")
		sendJsonResponse(w, response)
		return
	}

	fsm := &validate.ValidateFSM{}

	valid, results := fsm.Validate(r.Context(), reqEmail.Email)

	response := NewResponse(valid, results)
	sendJsonResponse(w, response)
}
