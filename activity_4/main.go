package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"Ephraim_PaxosLab/paxos"
)

// Replace these with your Minikube IP and NodePort for each replica
var acceptorURLs = []string{
	"http://paxos-service:8080",
}

type RequestBody struct {
	ProposalNumber int    `json:"ProposalNumber"`
	Value          string `json:"Value"`
}

func proposeHandler(w http.ResponseWriter, r *http.Request) {
	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body: %v\n", err)
		return
	}

	// 2-second timeout for each HTTP request
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Prepare phase
	promises := 0
	for _, url := range acceptorURLs {
		_, err := paxos.SendPrepareHTTPWithTimeout(ctx, url, body.ProposalNumber)
		if err != nil {
			fmt.Printf("Prepare failed for %s: %v\n", url, err)
			continue
		}
		promises++
	}

	if promises <= len(acceptorURLs)/2 {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "Consensus not reached: not enough promises\n")
		return
	}

	// Accept phase
	accepted := 0
	for _, url := range acceptorURLs {
		_, err := paxos.SendAcceptHTTPWithTimeout(ctx, url, body.ProposalNumber, body.Value)
		if err != nil {
			fmt.Printf("Accept failed for %s: %v\n", url, err)
			continue
		}
		accepted++
	}

	if accepted > len(acceptorURLs)/2 {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Consensus reached: %s\n", body.Value)
	} else {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "Consensus not reached: not enough accepts\n")
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Paxos Web Service (Fault Tolerant) running on port", port)
	http.HandleFunc("/propose", proposeHandler)
	http.ListenAndServe(":"+port, nil)
}
