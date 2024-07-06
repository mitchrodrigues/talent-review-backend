package openai

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/golly-go/golly"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOpenAIClient_request(t *testing.T) {
	gctx := golly.NewContext(context.Background())

	// Create an instance of our test object
	mockClient := new(MockHTTPClient)
	client := OpenAIClient{HTTPClient: mockClient}

	// Create a mock response
	responseBody := `{"some": "response"}`

	r := io.NopCloser(bytes.NewReader([]byte(responseBody)))
	mockResponse := &http.Response{StatusCode: http.StatusOK, Body: r}

	// Setup expectations
	mockClient.On("Do", mock.Anything).Return(mockResponse, nil)

	// Create a mock context and payload

	payload := CompletionPayload{}

	// Call the method
	response, err := client.request(gctx, "", payload)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, []byte(responseBody), response)

	mockClient.AssertExpectations(t)
}
