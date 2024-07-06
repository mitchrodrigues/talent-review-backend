package openai

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewDefaultContext tests the NewDefaultContext function
func TestNewDefaultContext(t *testing.T) {
	who := "TestWho"
	action := "TestAction"
	content := "TestContent"

	result := NewDefaultContext(who, action, content)

	assert.Equal(t, RoleSystem, result.Role, "Role should be set to RoleSystem")
	assert.Equal(t, who, result.Who, "Who should match the input")
	assert.Equal(t, action, result.Action, "Action should match the input")
	assert.Equal(t, content, result.Content, "Content should match the input")
}

// TestNewDefaultContextWithRole tests the NewDefaultContextWithRole function
func TestNewDefaultContextWithRole(t *testing.T) {
	role := MessageRole("customRole")
	who := "TestWho"
	action := "TestAction"
	content := "TestContent"

	result := NewDefaultContextWithRole(role, who, action, content)

	assert.Equal(t, role, result.Role, "Role should match the input role")
	assert.Equal(t, who, result.Who, "Who should match the input")
	assert.Equal(t, action, result.Action, "Action should match the input")
	assert.Equal(t, content, result.Content, "Content should match the input")
}

func TestDefaultAIContext_String(t *testing.T) {
	testCases := []struct {
		name     string
		context  DefaultAIContext
		expected string
	}{
		{
			name:     "Empty Context",
			context:  DefaultAIContext{},
			expected: "",
		},
		{
			name:     "Only Who",
			context:  DefaultAIContext{Who: "TestWho"},
			expected: "",
		},
		{
			name:     "Who and Action",
			context:  DefaultAIContext{Who: "TestWho", Action: "Testing"},
			expected: "",
		},
		{
			name:     "Content Only",
			context:  DefaultAIContext{Content: "Hello"},
			expected: "Hello",
		},
		{
			name:     "Full Context",
			context:  DefaultAIContext{Who: "TestWho", Action: "Testing", Content: "Hello"},
			expected: "TestWho Testing Hello",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.context.String())
		})
	}
}

// Test cases for Message() method
func TestDefaultAIContext_Message(t *testing.T) {
	testCases := []struct {
		name     string
		context  DefaultAIContext
		expected Message
	}{
		{
			name:     "Empty Role",
			context:  DefaultAIContext{Content: "Hello"},
			expected: Message{Role: RoleSystem, Content: "Hello"},
		},
		{
			name:     "Specified Role",
			context:  DefaultAIContext{Role: RoleUser, Content: "Hello"},
			expected: Message{Role: RoleUser, Content: "Hello"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.context.Message())
		})
	}
}
