package openai

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golly-go/golly"
	"github.com/sirupsen/logrus"
)

// Pretty sure i hate this but will do this later

type CompletionPrompt interface {
	Prompt

	AIPretendsToBe() string
	AIScenerioContext() string

	DefaultAIContext() AIContexts
}

type CompletionPromptBase struct {
	PromptBase
}

func (CompletionPromptBase) AIPretendsToBe() string {
	return config.AIPretendsToBe
}

func (CompletionPromptBase) AIScenerioContext() string {
	return config.AIScenarioContext
}

func (p CompletionPromptBase) DefaultAIContext() AIContexts {
	ret := AIContexts{}

	if str := p.AIPretendsToBe(); str != "" {
		ret = append(ret, NewDefaultSystemContext(str))
	}

	if str := p.AIScenerioContext(); str != "" {
		ret = append(ret, NewDefaultSystemContext(str))
	}

	return ret
}

// Define a structure to match the JSON payload structure
type CompletionPayload struct {
	Model    AIModel         `json:"model"`
	Format   *ResponseFormat `json:"response_format,omitempty"`
	Messages []Message       `json:"messages"`

	Seed *int64 `json:"seed,omitempty"`

	Temperature float64 `json:"temperature,omitempty"`
	N           float64 `json:"n"`

	FrequencyPenalty *float64           `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]float64 `json:"logit_bias,omitempty"`
}

type CompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`

	SystemFingerprint string `json:"system_fingerprint"`

	Choices []Choice       `json:"choices"`
	Usage   Usage          `json:"usage"`
	Model   AIModel        `json:"model"`
	Error   *ErrorResponse `json:"error,omitempty"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Message struct to hold each message
type Message struct {
	Role    MessageRole `json:"role"`
	Content string      `json:"content"`
}

func (m Message) Message() Message {
	return m
}

type Messages []Message

func (m Messages) Append(messages ...Message) Messages {
	return append(m, messages...)
}

func NewMessageArray(messages ...Message) Messages {
	return messages
}

// SendRequest sends a request to the OpenAI API and returns the response
func (oai OpenAIClient) Completions(ctx golly.Context, payload CompletionPayload) (CompletionResponse, error) {
	defer func(start time.Time) {
		ctx.Logger().
			WithFields(logrus.Fields{"model": payload.Model}).
			Infof("[%s] OpenAI Request took %s", payload.Model, time.Since(start))
	}(time.Now())

	if payload.Model == AutoModel {
		payload.Model = FastModel
	}

	body, err := oai.request(ctx, openAICompletionURL, payload)
	if err != nil {
		return CompletionResponse{}, err
	}

	aiResp := CompletionResponse{}
	if err := json.Unmarshal(body, &aiResp); err != nil {
		return aiResp, fmt.Errorf("error unmarshalling response: %w", err)
	}

	if aiResp.Error != nil {
		return aiResp, aiResp.Error
	}

	return aiResp, nil
}
