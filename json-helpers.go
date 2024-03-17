package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type JsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	// limit the size of the uploaded json file
	maxByte := 1048576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxByte))
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body should have only one single JSON value")
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
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
		return err
	}
	return nil
}

func ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}
	var payload JsonResponse
	payload.Error = true
	payload.Message = err.Error()
	return WriteJSON(w, statusCode, payload)
}

// /////////////////////////////////////// USAGE /////////////////////////////////
// example:
/*
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := ReadJSON(w,r, &requestPayload)
	if err != nil {
		ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	payload := JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}
	WriteJSON(w, http.StatusAccepted, payload)
}
*/
