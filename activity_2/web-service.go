package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"Ephraim_PaxosLab/paxos"
)

var (
	acceptors = []*paxos.Acceptor{
		&paxos.Acceptor{},
		&paxos.Acceptor{},
		&paxos.Acceptor{},
	}
	// acceptors = []string{
	// "http://localhost:8080",
	// "http://localhost:8081",
	// "http://localhost:8082",
	// }

	mu sync.Mutex
)

type RequestBody struct {
	ProposalNumber int    `json:"ProposalNumber"`
	Value          string `json:"Value"`
}

func proposeHandler(w http.ResponseWriter, r *http.Request) {
	var body RequestBody

	// Decode request body safely
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body: %v\n", err)
		return
	}

	// Create proposer
	proposer := paxos.Proposer{
		ProposalNumber: body.ProposalNumber,
		Value:          body.Value,
	}

	// Paxos must be synchronized
	mu.Lock()
	value := proposer.Propose(body.Value, acceptors)
	mu.Unlock()

	// Response
	if value != nil {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Consensus reached: %s\n", value)
	} else {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "Consensus not reached\n")
	}
}

func main() {
	fmt.Println("Paxos Web Service Running on port 8080...")
	http.HandleFunc("/propose", proposeHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)

}
