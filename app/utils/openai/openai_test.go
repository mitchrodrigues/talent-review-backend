package openai

import (
	"context"
	"testing"

	"github.com/golly-go/golly"
	"github.com/stretchr/testify/assert"
)

func TestLLM(t *testing.T) {
	// Create a mock context
	ctx := golly.NewContext(context.TODO())

	// First call to LLM - should create a new client
	client1 := LLM(ctx)
	assert.NotNil(t, client1, "Expected non-nil client")

	// Retrieve the client directly from the context
	retrievedClient, ok := ctx.Get(openAIContextKey)
	assert.True(t, ok, "Expected to find a client in the context")
	assert.Equal(t, client1, retrievedClient, "Expected the client to be the same as returned by LLM")

	// Second call to LLM - should return the same client
	client2 := LLM(ctx)
	assert.Equal(t, client1, client2, "Expected LLM to return the same client on subsequent calls")
}
