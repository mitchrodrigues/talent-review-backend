package openai

import (
	"encoding/json"
	"time"

	"github.com/golly-go/golly"
	"github.com/golly-go/golly/errors"
)

type EmbeddingResponse struct {
	Data   []EmbeddingData `json:"data"`
	Model  string          `json:"model"`
	Object string          `json:"object"`
	Usage  Usage           `json:"usage"`

	Error *ErrorResponse `json:"error,omitempty"`
}

type EmbeddingData struct {
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
	Object    string    `json:"object"`
}

func (oai OpenAIClient) Embeddings(ctx golly.Context, text string) (EmbeddingResponse, error) {

	defer func(start time.Time) {
		ctx.Logger().
			Infof("Embedding Request took %s", time.Since(start))
	}(time.Now())

	response := EmbeddingResponse{}

	params := map[string]any{"model": EmbeddingModel, "input": text}

	body, err := oai.request(ctx, openAIEmbeddingURL, params)
	if err != nil {
		return response, errors.WrapGeneric(err)
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return response, errors.WrapGeneric(err)
	}

	if response.Error != nil {
		return response, response.Error
	}

	return response, nil
}
