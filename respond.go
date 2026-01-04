package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type errResp struct {
	Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, status int, msg string) {
	if status > 499 {
		log.Println("Responding with 5XX error: ", msg)
	}

	respondWithJson(w, status, errResp{
		Error: msg,
	})

}
func respondWithJson[T any](w http.ResponseWriter, status int, payload T) error {
	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("failed to marshal JSON response with payload %v: error: %v", payload, err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
	return nil
}
