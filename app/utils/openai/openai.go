package openai

import (
	"github.com/golly-go/golly"
)

type AIModel string
type MessageRole string

const (
	LargeTokenModel AIModel = "gpt-4o"
	StandardModel   AIModel = "gpt-4o"
	FastModel       AIModel = "gpt-4o"
	TurboModel      AIModel = "gpt-4o"

	EmbeddingModel AIModel = "text-embedding-ada-002"
	WhisperModel   AIModel = "whipser-1"
	TTSModel       AIModel = "tts-1"
	ImageModel     AIModel = "dall-e-3"

	AutoModel AIModel = ""

	RoleUser   MessageRole = "user"
	RoleSystem MessageRole = "system"
	RoleAI     MessageRole = "assistant"

	openAICompletionURL = "https://api.openai.com/v1/chat/completions"
	openAIEmbeddingURL  = "https://api.openai.com/v1/embeddings"
	openAITTSURL        = "https://api.openai.com/v1/audio/speech"

	openAIContextKey golly.ContextKeyT = "openaiClient"
)

var (
	EmbeddingZero = [][]float64{
		{0.0},
	}

	FormatJSON = ResponseFormat{Type: "json_object"}

	config Config
)

type Config struct {
	// This will be picked up from ENV OPENAI_TOKEN
	// this is just here to seed a default.
	Token string

	AIPretendsToBe    string
	AIScenarioContext string
}

type ResponseFormat struct {
	Type string `json:"type"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ErrorResponse struct {
	Message string  `json:"message"`
	Type    string  `json:"type"`
	Param   *string `json:"param"`
	Code    string  `json:"code"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}

func SystemMessage(content string) Message {
	return Message{Role: RoleSystem, Content: content}
}

func UserMessage(content string) Message {
	return Message{Role: RoleUser, Content: content}
}

// Maybe hmm
func LLM(gctx golly.Context) *OpenAIClient {
	if client, ok := gctx.Get(openAIContextKey); ok {
		return client.(*OpenAIClient)
	}

	client := NewClient(gctx)
	gctx.Set(openAIContextKey, client)

	return client
}

func NewClient(ctx golly.Context) *OpenAIClient {
	var client HTTPClient = &CustomHTTPClient{}
	if golly.Env().IsTest() {
		client = &MockHTTPClient{}
	}

	return &OpenAIClient{
		HTTPClient: client,
		Token:      config.Token,
	}
}

func Initailizer(c Config) error {
	config = c
	return nil
}
