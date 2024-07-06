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

type TestPrompt struct {
	CompletionPromptBase `ai:"-"`

	rules    []string
	scenario []string
	context  AIContexts

	AIField string `json:"aiField"`
}

func (tp TestPrompt) Rules(ctx golly.Context) []string {
	// Return rules based on the test prompt.
	return tp.rules
}

func (tp TestPrompt) Scenario(ctx golly.Context) []string {
	// Return scenario based on the test prompt.
	return tp.scenario
}

func (tp TestPrompt) Context(ctx golly.Context) AIContexts {
	// Return context based on the test prompt.
	return tp.context
}

func TestFields(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected map[string]interface{}
	}{
		{
			name: "Simple Struct",
			input: struct {
				Field1 string `json:"field1" ai:"field2"`
				Field2 int    `ai:"field2"`
			}{Field1: "value1", Field2: 2},
			expected: map[string]interface{}{
				"field1": "field2",
				"Field2": "field2",
			},
		},
		{
			name: "Nested Struct",
			input: struct {
				SimpleField struct {
					Field1 string `json:"field1" ai:"field2"`
					Field2 int    `ai:"field2"`
				} `json:"simple"`
			}{SimpleField: struct {
				Field1 string `json:"field1" ai:"field2"`
				Field2 int    `ai:"field2"`
			}{
				Field1: "value1",
				Field2: 2,
			},
			},

			expected: map[string]interface{}{
				"simple": map[string]interface{}{
					"field1": "field2",
					"Field2": "field2",
				},
			},
		},
		{
			name: "Struct With Tag Omit",
			input: struct {
				Field1 string `json:"field1"`
				Field2 int    `json:"field2" ai:"-"`
			}{Field1: "value1", Field2: 2},
			expected: map[string]interface{}{
				"field1": "string",
			},
		},
		{
			name: "Test Type with Array",
			input: struct {
				Field1 []struct {
					Field1 string `json:"field1"`
				} `json:"field1"`
			}{},
			expected: map[string]interface{}{
				"field1": []interface{}{
					map[string]interface{}{
						"field1": "string",
					},
				},
			},
		},
		// Add more test cases as needed
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := Fields(tc.input)

			assert.EqualValues(t, tc.expected, result)
		})
	}
}

func TestBindings(t *testing.T) {
	tests := []struct {
		name         string
		prompt       Prompt
		instructions []string
		expected     map[string]interface{}
	}{
		{
			name: "Test Case 1",
			prompt: TestPrompt{
				rules: []string{"rule 1"},
			},
			instructions: []string{"Instruction 1", "Instruction 2"},
			expected: map[string]interface{}{
				"fields":   "{\n\t\"aiField\": \"string\"\n}",
				"rules":    []string{"rule 1"},
				"scenario": []string{"Instruction 1", "Instruction 2"},
			},
		},
		// Additional test cases.
	}

	ctx := golly.NewContext(context.Background())

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := Bindings(ctx, tc.prompt, tc.instructions...)

			// Compare each key in the result map with the expected map.

			assert.EqualValues(t, tc.expected, result)
		})
	}
}

func TestCompletion(t *testing.T) {

	// Define your test cases
	tests := []struct {
		name     string
		prompt   CompletionPrompt
		expected CompletionPrompt

		instructions []string
		mockResponse *http.Response
		expectError  bool
	}{
		{
			name:   "Test Case 1",
			prompt: &TestPrompt{scenario: []string{"Instruction 1", "Instruction 2"}},
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewReader([]byte(`{
					"choices": [
					   {
						   "index": 0,
						   "message": {
							   "role": "assistant",
							   "content": "{\"aiField\": \"string\"}"
						   },
						   "finish_reason": "stop"
					   }
				   ]}`,
				))),
			},
			expected:    &TestPrompt{AIField: "string", scenario: []string{"Instruction 1", "Instruction 2"}},
			expectError: false,
		},
		// Additional test cases.
	}

	gctx := golly.NewContext(context.Background())

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := new(MockHTTPClient)

			client := OpenAIClient{HTTPClient: mockClient}

			// Setup expectations for the HTTP client
			mockClient.On("Do", mock.Anything).Return(tc.mockResponse, nil)

			// Set the mock client to LLM or its underlying HTTP client
			// ...

			result, err := Completion(gctx, &client, tc.prompt)

			if tc.expectError {
				// Check if error is expected
				assert.Error(t, err, "Error expectation mismatch")
			} else {
				assert.NoError(t, err)
			}

			// Compare the result with the expected prompt
			assert.Equal(t, tc.expected, result, "Unexpected result from Completion")
		})
	}
}
