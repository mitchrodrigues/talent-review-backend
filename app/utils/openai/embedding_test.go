package openai

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/golly-go/golly"
	"github.com/stretchr/testify/mock"
)

func TestOpenAIClient_Embeddings(t *testing.T) {
	testCases := []struct {
		name           string
		inputText      string
		mockResponse   *http.Response
		expectError    bool
		expectedOutput EmbeddingResponse
	}{
		{
			name:      "Valid Embedding",
			inputText: "example text",
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewReader([]byte(`{
                    "data": [{"embedding": [-0.006929283495992422, -0.005336422007530928], "index": 0, "object": "embedding"}],
                    "model": "text-embedding-ada-002",
                    "object": "list",
                    "usage": {"prompt_tokens": 5, "total_tokens": 5}
                }`))),
			},
			expectError: false,
			expectedOutput: EmbeddingResponse{
				Data: []EmbeddingData{
					{
						Embedding: []float64{-0.006929283495992422, -0.005336422007530928},
						Index:     0,
						Object:    "embedding",
					},
				},
				Model:  "text-embedding-ada-002",
				Object: "list",
				Usage:  Usage{PromptTokens: 5, TotalTokens: 5},
			},
		},
	}

	gctx := golly.NewContext(context.Background())

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := new(MockHTTPClient)
			client := OpenAIClient{HTTPClient: mockClient}

			// Setup expectations for the HTTP client
			mockClient.On("Do", mock.Anything).Return(tc.mockResponse, nil)

			result, err := client.Embeddings(gctx, tc.inputText)

			// Check if error is expected
			if (err != nil) != tc.expectError {
				t.Errorf("Test %s: Expected error: %v, got: %v", tc.name, tc.expectError, err)
			}

			// Check if the result matches the expected output
			if !reflect.DeepEqual(result, tc.expectedOutput) {
				t.Errorf("Test %s: Expected output: %+v, got: %+v", tc.name, tc.expectedOutput, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
