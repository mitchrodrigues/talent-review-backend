package openai

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/golly-go/golly"
	"github.com/golly-go/golly/errors"
	"github.com/golly-go/golly/utils"
)

func Bindings(ctx golly.Context, prompt Prompt, instructions ...string) map[string]interface{} {
	// Create an empty map to hold our results.
	result := make(map[string]interface{})

	result["fields"] = AIFieldsString(prompt)

	rules := prompt.Rules(ctx)
	if len(rules) > 0 {
		result["rules"] = rules
	}

	result["scenario"] = append(prompt.Scenario(ctx), instructions...)

	return result
}

func AIFieldsString(prompt Prompt) string {
	// Handle field data.
	fieldData := Fields(prompt)

	var field interface{}

	if reflect.TypeOf(prompt).Kind() == reflect.Slice {
		field = []interface{}{fieldData}
	} else {
		field = fieldData
	}

	b, _ := json.MarshalIndent(field, "", "\t")

	return string(b)
}

type CompletionOptions struct {
	Model       AIModel
	Format      *ResponseFormat
	Temperature float64

	LogitBias map[string]float64

	FrequencyPenalty float64
	PresencePenalty  float64
	N                float64
}

func buildCompletionPayload[T CompletionPrompt](gctx golly.Context, prompt T, aiCtxs ...AIContext) (CompletionPayload, error) {
	var frequencyPenalty *float64 = nil
	var temperature float64

	pstring, err := prompt.PromptString(gctx, prompt)
	if err != nil {
		return CompletionPayload{}, errors.WrapGeneric(err)
	}

	aiContexts := AIContexts{}

	if fp := prompt.FrequencyPenalty(); fp != -1 {
		frequencyPenalty = &fp
	}

	if temperature = prompt.Temperature(); temperature == -1 {
		temperature = 0.4
	}

	gctx.Logger().Debugf("Prompt: \n%s\n\n", pstring)

	messages := aiContexts.
		Append(prompt.DefaultAIContext()...).
		Append(aiCtxs...).
		Append(golly.Compact(aiContextForPrompts(gctx, prompt))...).
		Messages()

	return CompletionPayload{
		Model:            prompt.Model(),
		Format:           &FormatJSON,
		Temperature:      temperature,
		N:                prompt.N(),
		LogitBias:        prompt.LogitBias(),
		FrequencyPenalty: frequencyPenalty,
		Messages: messages.Append(Message{
			Role:    RoleUser,
			Content: pstring,
		}),
	}, nil
}

func Completion[T CompletionPrompt](gctx golly.Context, llm *OpenAIClient, prompt T, aiCtxs ...AIContext) (T, error) {
	payload, err := buildCompletionPayload(gctx, prompt, aiCtxs...)
	if err != nil {
		return prompt, errors.WrapGeneric(err)
	}

	result, err := llm.Completions(gctx, payload)

	if err != nil {
		return prompt, errors.WrapGeneric(err)
	}

	if len(result.Choices) == 0 {
		return prompt, errors.WrapGeneric(fmt.Errorf("no choices returned"))
	}

	gctx.Logger().Debugf("Response %s", result.Choices[0].Message.Content)

	if err := marshalPrompt(result.Choices[0].Message.Content, prompt); err != nil {
		return prompt, errors.WrapGeneric(err)
	}

	return prompt, nil
}

func marshalPrompt(response string, prompt interface{}) error {
	if reflect.ValueOf(prompt).Kind() == reflect.Ptr {
		return json.Unmarshal([]byte(response), prompt)
	}
	return json.Unmarshal([]byte(response), &prompt)
}

func Fields(prompt any) map[string]interface{} {
	tpe := reflect.ValueOf(prompt)

	// Simplify how we handle slice or pointer to get to the actual type
	if tpe.Kind() == reflect.Slice {
		tpe = reflect.New(tpe.Type().Elem()).Elem()
	} else if tpe.Kind() == reflect.Ptr {
		tpe = tpe.Elem()
	}

	ret := make(map[string]interface{})

	for pos := 0; pos < tpe.NumField(); pos++ {
		field := tpe.Type().Field(pos)
		val := tpe.Field(pos)

		name := field.Name

		aiTag := field.Tag.Get("ai")
		jsonTag := field.Tag.Get("json")

		if !field.IsExported() || aiTag == "-" || jsonTag == "-" {
			continue
		}

		if jsonTag != "" {
			name = jsonTag
		}

		switch val.Kind() {
		case reflect.Struct:
			ret[name] = Fields(val.Interface())
		case reflect.Slice:
			t := val.Type().Elem()
			if t.Kind() == reflect.Struct {
				if hsh := Fields(reflect.New(t).Elem().Interface()); len(hsh) > 0 {
					ret[name] = []interface{}{hsh}
				}
				continue
			}
			ret[name] = []interface{}{utils.GetType(val.Interface())}
		default:
			if aiTag != "" {
				ret[name] = aiTag
			} else {
				ret[name] = utils.GetType(val.Interface())
			}
		}
	}
	return ret
}
