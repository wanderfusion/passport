package httpclient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// MakePOSTRequest makes a POST request and returns the response as a byte array.
func MakePOSTRequest[T any](client *http.Client, baseURL string, endpoint string, payload T, customHeaders map[string]string, customQueryParams map[string]string) ([]byte, error) {
	// Marshal the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullURL := baseURL + endpoint

	// Create the request
	req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	// Add custom headers if provided
	for key, value := range customHeaders {
		req.Header.Set(key, value)
	}

	// Add custom query parameters if provided
	q := req.URL.Query()
	for key, value := range customQueryParams {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}
