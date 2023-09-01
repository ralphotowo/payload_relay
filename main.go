package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// Define the port to listen on and the endpoint to forward payloads to.
	port := ":8080"
	endpointURL := "https://eo84zdkagfp0p87.m.pipedream.net"

	// Create an HTTP server to listen on the specified port.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if the request path is "/healthcheck" and do not forward it.
		if r.URL.Path == "/healthcheck" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Request to /healthcheck is not forwarded"))
			return
		}

		// Read the payload from the incoming request.
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Forward the payload to the specified endpoint.
		resp, err := forwardPayload(endpointURL, payload)
		if err != nil {
			http.Error(w, "Failed to forward payload", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Return the response from the endpoint as the response to the incoming request.
		responsePayload, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Failed to read response body from endpoint", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		w.Write(responsePayload)
	})

	// Start the HTTP server.
	fmt.Printf("Listening on port %s...\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func forwardPayload(endpointURL string, payload []byte) (*http.Response, error) {
	resp, err := http.Post(endpointURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	return resp, nil
}
