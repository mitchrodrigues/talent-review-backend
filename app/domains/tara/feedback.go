package tara

import (
	"github.com/golly-go/golly"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/openai"
)

type SummarizeFeedbackInput struct {
	Strengths          string
	Opportunities      string
	AdditionalComments string
}

type ActionItem struct {
	Item string `json:"item"  ai:"string action items for the manager"`
}

type ActionItems []ActionItem

func (ai ActionItems) Values() []string {
	return golly.Map(ai, func(item ActionItem) string {
		return item.Item
	})
}

type SummarizeFeedbackPrompt struct {
	openai.CompletionPromptBase `json:"-"`
	SummarizeFeedbackInput      `json:"-"`

	Summary     string      `json:"summary" ai:"string summary of the feedback"`
	ActionItems ActionItems `json:"action_items"`
}

func (SummarizeFeedbackPrompt) Rules(gctx golly.Context) []string {
	return []string{
		"The summary should be concise and cover key points of the feedback.",
		"Identify and highlight any action items for the manager to review.",
		"Use simple,concise, and professional wording.",
		"If their are examples in the feedback be sure to include them in the summary",
		"Do not editorialize the feedback.",
		"Do not omit details from the feedback",
	}
}

func (prompt SummarizeFeedbackPrompt) Scenario(gctx golly.Context) []string {
	return []string{
		"You are given the following feedback about an employee:",
		"Strengths: " + prompt.Strengths,
		"Opportunities: " + prompt.Opportunities,
		"Additional Comments: " + prompt.AdditionalComments,
		"Summarize the feedback clearly.",
		"Call out any action items that may exist for the manager to review.",
		"Ensure that the summary is simple, concise, and uses professional wording without editorializing.",
	}
}

func NewSummaryFeedbackPrompt(input SummarizeFeedbackInput) *SummarizeFeedbackPrompt {
	return &SummarizeFeedbackPrompt{
		SummarizeFeedbackInput: input,
	}
}
