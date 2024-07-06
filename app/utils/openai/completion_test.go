package openai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/golly-go/golly"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOpenAIClient_Completions(t *testing.T) {
	type testCase struct {
		name           string
		payload        CompletionPayload
		mockResponse   *http.Response
		expectedError  error
		expectedOutput []byte // Or any other type that matches your expected output
	}

	tests := []testCase{
		{
			name:    "No Error",
			payload: CompletionPayload{ /* fill with test-specific data */ },
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewReader([]byte(`{
					 "choices": [
						{
							"index": 0,
							"message": {
								"role": "assistant",
								"content": "Thank you for the expression of affection!"
							},
							"finish_reason": "stop"
						}
					]}`,
				))),
			},
			expectedOutput: []byte(`expected output for test case 1`),
		},
		{
			name:    "API Error",
			payload: CompletionPayload{ /* fill with test-specific data */ },
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewReader([]byte(`{
						"error": {
							"message": "The model gpt-4.0 does not exist",
							"type": "invalid_request_error",
							"param": null,
							"code": "model_not_found"
						}
					}`,
				))),
			},
			expectedError:  fmt.Errorf("The model gpt-4.0 does not exist"),
			expectedOutput: []byte(`expected output for test case 1`),
		},

		// Add more test cases as needed
	}

	gctx := golly.NewContext(context.Background())

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := new(MockHTTPClient)
			client := OpenAIClient{HTTPClient: mockClient}

			// Setup expectations for the HTTP client
			mockClient.On("Do", mock.Anything).Return(tc.mockResponse, nil)

			// Call the method
			aiResponse, err := client.Completions(gctx, tc.payload)

			// Check if error is expected

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Len(t, aiResponse.Choices, 1)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
