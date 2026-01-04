package main

import (
	"log"
	"net/http"
)

type resp struct {
	Status string `json:"status"` // w this tag -> instead of "Status" it will be json format "status"
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	err := respondWithJson(w, http.StatusOK, resp{Status: "OK server good 🧃 :)"})
	if err != nil {
		log.Printf("Error writing JSON response: %v", err)
	}
}

func handlerErr(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 400, "Something went wrong")
}
