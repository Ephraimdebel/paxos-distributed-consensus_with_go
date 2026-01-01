package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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

func sendPrepareHTTP(url string, proposalNumber int) (Promise, error) {
	body, _ := json.Marshal(Prepare{ProposalNumber: proposalNumber})
	resp, err := http.Post(url+"/prepare", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return Promise{}, err
	}
	defer resp.Body.Close()

	var promise Promise
	json.NewDecoder(resp.Body).Decode(&promise)
	return promise, nil
}

func sendAcceptHTTP(url string, proposalNumber int, value string) (Accepted, error) {
	body, _ := json.Marshal(Accept{ProposalNumber: proposalNumber, Value: value})
	resp, err := http.Post(url+"/accept", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return Accepted{}, err
	}
	defer resp.Body.Close()

	var ack Accepted
	json.NewDecoder(resp.Body).Decode(&ack)
	return ack, nil
}

func main() {
    acceptorURLs := []string{
        "http://localhost:8000",
        "http://localhost:8001",
        "http://localhost:8002",
        "http://localhost:8003",
        "http://localhost:8004",
    }

    proposalNumber := 10 
    value := "Distributed Systems"

    fmt.Println("===== PHASE 1: PREPARE =====")

    promises := 0
    alive := 0

    for _, url := range acceptorURLs {
        promise, err := sendPrepareHTTP(url, proposalNumber)
        if err != nil {
            fmt.Println(" Acceptor unreachable:", url)
            continue
        }

        alive++

        if promise.ProposalNumber == proposalNumber {
            fmt.Println("✔️ Promise from:", url)
            promises++
        } else {
            fmt.Println("Rejected by:", url)
        }
    }

    fmt.Printf("Result: %d promises from %d alive nodes (total=%d)\n",
        promises, alive, len(acceptorURLs))

    // Paxos rule: must be > N/2 of TOTAL nodes
    if promises <= len(acceptorURLs)/2 {
        fmt.Println(" Consensus not reached: not enough promises")
        return
    }

    fmt.Println("\n===== PHASE 2: ACCEPT =====")

    accepted := 0

    for _, url := range acceptorURLs {
        ack, err := sendAcceptHTTP(url, proposalNumber, value)
        if err != nil {
            fmt.Println(" Acceptor unreachable:", url)
            continue
        }

        if ack.ProposalNumber == proposalNumber {
            fmt.Println("✔️ Accepted by:", url)
            accepted++
        } else {
            fmt.Println(" Accept rejected by:", url)
        }
    }

    fmt.Printf("Result: %d accepts (total=%d)\n",
        accepted, len(acceptorURLs))

    if accepted > len(acceptorURLs)/2 {
        fmt.Printf(" Consensus reached on value: %s\n", value)
    } else {
        fmt.Println(" Consensus not reached: not enough accepts")
    }
}
