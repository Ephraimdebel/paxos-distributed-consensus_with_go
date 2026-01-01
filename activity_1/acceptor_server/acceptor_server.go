package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
)

type Prepare struct {
	ProposalNumber int `json:"proposal_number"`
}

type Promise struct {
	ProposalNumber int         `json:"proposal_number"`
	AcceptedValue  interface{} `json:"accepted_value"`
}

type Accept struct {
	ProposalNumber int         `json:"proposal_number"`
	Value          interface{} `json:"value"`
}

type Accepted struct {
	ProposalNumber int         `json:"proposal_number"`
	Value          interface{} `json:"value"`
}

type Acceptor struct {
	mu             sync.Mutex
	promisedNumber int
	acceptedNumber int
	acceptedValue  interface{}
}

var a Acceptor

func prepareHandler(w http.ResponseWriter, r *http.Request) {
	var p Prepare
	json.NewDecoder(r.Body).Decode(&p)

	a.mu.Lock()
	defer a.mu.Unlock()

	resp := Promise{}
	if p.ProposalNumber > a.promisedNumber {
		a.promisedNumber = p.ProposalNumber
		resp = Promise{
			ProposalNumber: p.ProposalNumber,
			AcceptedValue:  a.acceptedValue,
		}
	}

	json.NewEncoder(w).Encode(resp)
}

func acceptHandler(w http.ResponseWriter, r *http.Request) {
	var ac Accept
	json.NewDecoder(r.Body).Decode(&ac)

	a.mu.Lock()
	defer a.mu.Unlock()

	resp := Accepted{}
	if ac.ProposalNumber >= a.promisedNumber {
		a.promisedNumber = ac.ProposalNumber
		a.acceptedNumber = ac.ProposalNumber
		a.acceptedValue = ac.Value

		resp = Accepted{
			ProposalNumber: ac.ProposalNumber,
			Value:          ac.Value,
		}
	}

	json.NewEncoder(w).Encode(resp)
}

func main() {
	port := "8000" // default port
	if len(os.Args) > 1 {
		port = os.Args[1] // allow overriding port
	}

	http.HandleFunc("/prepare", prepareHandler)
	http.HandleFunc("/accept", acceptHandler)

	fmt.Println("Acceptor running on port", port)
	http.ListenAndServe(":"+port, nil)
}
