package openai

import (
	"fmt"
	"strings"

	"github.com/osteele/liquid"
)

var (
	engine = liquid.NewEngine()
)

type Template string

func (b Template) Generate(bindings map[string]interface{}) (string, error) {
	out, err := engine.ParseAndRenderString(string(b), bindings)
	if err != nil {
		fmt.Printf("problem with prompt: %s", string(b))
		return "", err
	}

	return strings.ReplaceAll(strings.ReplaceAll(out, "\t", " "), "\n\n", "\n"), nil
}
