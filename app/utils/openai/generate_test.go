package openai

import (
	"context"
	"testing"

	"github.com/golly-go/golly"
	"github.com/stretchr/testify/assert"
)

// TestBuildCompletionPayloadTable tests buildCompletionPayload function with table-driven approach
func TestBuildCompletionPayloadTable(t *testing.T) {
	// Define your test cases
	testCases := []struct {
		name            string
		mockPrompt      TestPrompt
		expectedPayload CompletionPayload
		expectedError   error
	}{
		{
			name:       "Valid Prompt",
			mockPrompt: TestPrompt{},
			expectedPayload: CompletionPayload{
				Model:       FastModel,
				Format:      &FormatJSON,
				N:           1,
				Temperature: 0.4,
				Messages: Messages{
					{Role: RoleUser, Content: "\nScenario:\nResult must be valid JSON in the following format:\n{\n \"aiField\": \"string\"\n}"},
				}},
			expectedError: nil,
		},

		// Add more test cases as necessary...
	}

	gctx := golly.NewContext(context.Background())

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			aiCtxs := []AIContext{} // Mock AI contexts if necessary

			// Call the function
			payload, err := buildCompletionPayload(gctx, tc.mockPrompt, aiCtxs...)

			// Assert the expectations
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedPayload, payload)
		})
	}
}
