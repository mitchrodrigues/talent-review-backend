package openai

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/golly-go/golly"
	"github.com/golly-go/golly/errors"
)

type Bool string

func (b Bool) True() bool {
	return strings.EqualFold(string(b), "yes") || strings.EqualFold(string(b), "true")
}

func (b Bool) False() bool {
	return !b.True()
}

func (b Bool) Value(deflt ...bool) bool {
	if b == "" && len(deflt) > 0 {
		return deflt[0]
	}

	return b.True()
}

type Int string

func (b Int) Value(deflt ...int) int {
	i, err := strconv.Atoi(string(b))

	if err != nil && len(deflt) > 0 {
		return deflt[0]
	}

	return i
}

type Float string

func (b Float) Value(deflt ...float64) float64 {
	f, err := strconv.ParseFloat(string(b), 64)

	if err != nil && len(deflt) > 0 {
		return deflt[0]
	}

	return f
}

type StringSlice struct {
	Values []string `json:"values" ai:"string values"`
}

func (s StringSlice) Value() []string {
	return s.Values
}

func NewStringSlice(values ...string) StringSlice {
	return StringSlice{Values: values}
}

type EmbeddableText string

func (t EmbeddableText) String() string {
	return string(t)
}

func (t EmbeddableText) Vecors(gctx golly.Context) ([]float64, error) {
	response, err := LLM(gctx).Embeddings(gctx, string(t))

	if err != nil {
		return []float64{}, errors.WrapGeneric(err)
	}

	if len(response.Data) == 0 {
		return []float64{}, errors.WrapNotFound(fmt.Errorf("no embeddings found"))
	}

	return response.Data[0].Embedding, nil
}
