package openai

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golly-go/golly"
	"github.com/golly-go/golly/errors"
)

const (
	openaiImageURL = "https://api.openai.com/v1/images/generations"

	MasterImagePrompt Template = `
	{% if rules %}Rules:{% for rule in rules %}
	{{ rule }}{% unless forloop.last %}{% endunless %}{% endfor %}
	Scenario:{% endif %}
	{% if style %}Generate a detailed {{ style }}.{% endif %}{% for field in scenario %}
	{{ field }}{% unless forloop.last %}{% endunless %}{% endfor %}`
)

type ImagePrompt interface {
	Prompt

	ArtStyle() string
	Scenario(golly.Context) []string
}

type ImagePromptBase struct {
	PromptBase
}

func (ImagePromptBase) ArtStyle() string {
	return "pixel art"
}

func (ImagePromptBase) PromptString(ctx golly.Context, prompt Prompt) (string, error) {
	var artStyle string
	if p, ok := prompt.(ImagePrompt); ok {
		artStyle = p.ArtStyle()
	}

	return MasterImagePrompt.Generate(map[string]interface{}{
		"scenario": prompt.Scenario(ctx),
		"rules":    prompt.Rules(ctx),
		"style":    artStyle,
	})
}

type ImagePayload struct {
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"`

	Model  AIModel `json:"model"`
	Prompt string  `json:"prompt"`
	N      int     `json:"n"`

	Size     string `json:"size,omitempty"`
	Quantity string `json:"quality,omitempty"`
}

type ImageResponse struct {
	Data []struct {
		URL           string `json:"url"`
		RevisedPrompt string `json:"revised_prompt"`
	} `json:"data"`

	Error *ErrorResponse `json:"error,omitempty"`
}

func (oai OpenAIClient) Image(ctx golly.Context, payload ImagePayload) (ImageResponse, error) {
	var result ImageResponse

	defer func(start time.Time) {
		ctx.Logger().
			Infof("Image Request took %s", time.Since(start))
	}(time.Now())

	payload.Model = ImageModel

	b, err := oai.request(ctx, openaiImageURL, payload)
	if err != nil {
		return result, err
	}

	ctx.Logger().Debugf("Image Response: %s", string(b))

	if err := json.Unmarshal(b, &result); err != nil {
		return result, errors.WrapGeneric(err)
	}

	if result.Error != nil {
		return result, fmt.Errorf(result.Error.Message)
	}

	return result, nil

}

func Image[T ImagePrompt](gctx golly.Context, llm *OpenAIClient, prompt T) (ImageResponse, error) {

	pstring, err := prompt.PromptString(gctx, prompt)
	if err != nil {
		return ImageResponse{}, errors.WrapGeneric(err)
	}

	result, err := llm.Image(gctx, ImagePayload{
		Prompt: pstring,
		N:      int(prompt.N()),
	})

	if err != nil {
		return ImageResponse{}, errors.WrapGeneric(err)
	}

	if len(result.Data) == 0 {
		return ImageResponse{}, errors.WrapGeneric(fmt.Errorf("no choices returned"))
	}

	return result, nil
}
