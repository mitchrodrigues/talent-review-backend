package openai

import (
	"github.com/golly-go/golly"
)

// Sorta setting this up to be a bit dynamic prompt generator we can create and store prompts with actions and events later
const (
	MasterPrompt Template = `{% if rules %}Rules:{% for rule in rules %}
{{ rule }}{% unless forloop.last %}{% endunless %}{% endfor %}{% endif %}
Scenario:{% for field in scenario %}
{{ field }}{% unless forloop.last %}{% endunless %}{% endfor %}

Result must be valid JSON in the following format:
{{ fields }}`
)

type Prompt interface {
	Context(golly.Context) AIContexts

	PreviousPromptContexts(ctx golly.Context, prompt Prompt) AIContexts

	Scenario(golly.Context) []string
	Rules(golly.Context) []string

	PromptString(golly.Context, Prompt) (string, error)
	PromptToContexts(golly.Context) AIContexts

	Model() AIModel
	Temperature() float64
	LogitBias() map[string]float64
	FrequencyPenalty() float64
	PresencePenalty() float64
	ArtStyle() string

	N() float64
}

type PromptBase struct {
	PreviousPrompts []Prompt `json:"-" ai:"-"`
}

func (PromptBase) DefaultAIContext() AIContexts {
	return AIContexts{}
}

func (PromptBase) ArtStyle() string {
	return ""
}

func (pb PromptBase) Model() AIModel {
	return FastModel
}

func (pb PromptBase) Context(ctx golly.Context) AIContexts {
	return AIContexts{}
}

func (pb PromptBase) Temperature() float64 {
	return -1
}

func (pb PromptBase) LogitBias() map[string]float64 {
	return nil
}

func (pb PromptBase) FrequencyPenalty() float64 {
	return -1
}

func (pb PromptBase) PresencePenalty() float64 {
	return -1
}

func (pb PromptBase) N() float64 {
	return 1
}

func (pb PromptBase) Rules(ctx golly.Context) []string {
	return []string{}
}

func (pb *PromptBase) AddPreviousPrompts(prompts ...Prompt) {
	pb.PreviousPrompts = append(pb.PreviousPrompts, prompts...)
}

func (pb PromptBase) PreviousPromptContexts(ctx golly.Context, prompt Prompt) AIContexts {
	contexts := AIContexts{}

	for _, p := range pb.PreviousPrompts {
		contexts = append(contexts, p.PromptToContexts(ctx)...)
	}

	return contexts
}

func (pb PromptBase) PromptToContexts(ctx golly.Context) AIContexts {
	return AIContexts{}
}

func (pb PromptBase) PromptString(ctx golly.Context, prompt Prompt) (string, error) {
	return MasterPrompt.Generate(Bindings(ctx, prompt))
}

func aiContextForPrompts(ctx golly.Context, prompt Prompt) AIContexts {
	aictx := AIContexts{}

	aictx.
		Append(prompt.Context(ctx)...).
		Append(prompt.PreviousPromptContexts(ctx, prompt)...)

	return aictx
}
