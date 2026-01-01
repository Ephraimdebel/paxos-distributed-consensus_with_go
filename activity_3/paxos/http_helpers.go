package paxos

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Timeout and retry constants
const (
	RequestTimeout = 2 * time.Second
	MaxRetries     = 3
)

func SendPrepareHTTPWithTimeout(ctx context.Context, url string, proposalNumber int) (Promise, error) {
	var promise Promise

	for attempt := 1; attempt <= MaxRetries; attempt++ {
		reqBody := Prepare{ProposalNumber: proposalNumber}
		data, _ := json.Marshal(reqBody)

		req, err := http.NewRequestWithContext(ctx, "POST", url+"/prepare", bytes.NewReader(data))
		if err != nil {
			return promise, err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("Attempt %d: Failed to reach %s: %v\n", attempt, url, err)
			time.Sleep(500 * time.Millisecond)
			continue
		}

		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&promise); err != nil {
			return promise, err
		}

		return promise, nil
	}

	return promise, fmt.Errorf("failed to reach %s after %d retries", url, MaxRetries)
}

func SendAcceptHTTPWithTimeout(ctx context.Context, url string, proposalNumber int, value interface{}) (Accepted, error) {
	var ack Accepted

	for attempt := 1; attempt <= MaxRetries; attempt++ {
		reqBody := Accept{ProposalNumber: proposalNumber, Value: value}
		data, _ := json.Marshal(reqBody)

		req, err := http.NewRequestWithContext(ctx, "POST", url+"/accept", bytes.NewReader(data))
		if err != nil {
			return ack, err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("Attempt %d: Failed to reach %s: %v\n", attempt, url, err)
			time.Sleep(500 * time.Millisecond)
			continue
		}

		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&ack); err != nil {
			return ack, err
		}

		return ack, nil
	}

	return ack, fmt.Errorf("failed to reach %s after %d retries", url, MaxRetries)
}
