package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/golly-go/golly"
)

var (
	httpClient = &http.Client{}
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type CustomHTTPClient struct{}

func (c *CustomHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return httpClient.Do(req)
}

type OpenAIClient struct {
	Token      string
	HTTPClient HTTPClient

	Temperature float64
}

func (oai OpenAIClient) request(gctx golly.Context, url string, payload any) ([]byte, error) {
	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return []byte{}, fmt.Errorf("error marshalling the payload: %w", err)
	}

	gctx.Logger().Debugf("LLM Request: %s", string(payloadBytes))

	// Create a new request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return []byte{}, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+oai.Token)

	// Send the request using the same client to leverage HTTP Keep-Alive
	resp, err := oai.HTTPClient.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("error sending request: %w", err)
	}

	defer resp.Body.Close()

	// Read the response
	return io.ReadAll(resp.Body)
}

var _ HTTPClient = &CustomHTTPClient{}
