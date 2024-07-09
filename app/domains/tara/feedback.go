package tara

import (
	"fmt"

	"github.com/golly-go/golly"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/openai"
)

type SummarizeFeedbackInput struct {
	Strengths          string
	Opportunities      string
	AdditionalComments string
}

type ActionItem struct {
	Item string `json:"item"  ai:"string action item for the manager"`
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

	Summary string `json:"summary" ai:"string summary of the feedback"`
}

func (prompt SummarizeFeedbackPrompt) Context(gctx golly.Context) openai.AIContexts {
	return openai.AIContexts{
		openai.NewDefaultContentRole(openai.RoleUser,
			fmt.Sprintf(
				"Strengths: %s\nOpportunities: %s\nAdditional Comments: %s",
				prompt.SummarizeFeedbackInput.Strengths,
				prompt.SummarizeFeedbackInput.Opportunities,
				prompt.SummarizeFeedbackInput.AdditionalComments,
			),
		),
	}
}

func (SummarizeFeedbackPrompt) Rules(gctx golly.Context) []string {
	return []string{
		"The summary should be concise and cover key points of the feedback.",
		"Use simple, concise, and professional wording.",
		"If there are examples in the feedback, be sure to include them in the summary.",
		"Do not editorialize the feedback.",
		"Do not omit details from the feedback.",
	}
}

func (prompt SummarizeFeedbackPrompt) Scenario(gctx golly.Context) []string {
	return []string{
		"You are given feedback about an employee. Summarize the feedback clearly.",
		"Ensure that the summary is simple, concise, and uses professional wording without editorializing.",
	}
}

func (prompt SummarizeFeedbackPrompt) PromptToContexts(gctx golly.Context) openai.AIContexts {
	return append(prompt.Context(gctx),
		openai.NewDefaultContentRole(openai.RoleAI, prompt.Summary))
}

func NewSummaryFeedbackPrompt(input SummarizeFeedbackInput) *SummarizeFeedbackPrompt {
	return &SummarizeFeedbackPrompt{
		SummarizeFeedbackInput: input,
	}
}

type FollowUpItemsPrompt struct {
	openai.CompletionPromptBase `json:"-"`

	FollowUpItems ActionItems `json:"follow_up_items" ai:"string follow-up items for the manager"`
}

func (FollowUpItemsPrompt) Rules(gctx golly.Context) []string {
	return []string{
		"Identify and highlight specific follow-up items for the manager to review based on the provided feedback.",
		"Highlight any ambiguities, lack of details, or concerning issues that the manager should follow up on.",
		"Follow-up items should be specific tasks for the manager, such as clarifying ambiguous feedback, addressing lack of details, or taking action on concerning issues.",
		"Do not include career advice or general improvement suggestions for the employee.",
		"Use simple, concise, and professional wording.",
	}
}

func (prompt FollowUpItemsPrompt) Scenario(gctx golly.Context) []string {
	return []string{
		"You are given feedback about an employee. Based on this feedback, identify specific follow-up items for the manager to review.",
		"Highlight any ambiguities, lack of details, or concerning issues that the manager should address.",
		"Ensure that follow-up items are specific tasks for the manager and do not include career advice or general improvement suggestions for the employee.",
		"Example follow-up items:",
		"1. Clarify the specific areas in which the employee needs to improve communication skills.",
		"2. Follow up to get more details on reported delays in project delivery.",
		"3. Investigate the reasons behind the employee's inconsistent performance.",
		"4. Schedule a one-on-one meeting to discuss the employee's concerns about team collaboration.",
		"5. Ensure the employee receives specific examples of where they excel and where they need improvement.",
	}
}

func NewFollowUpItemsPrompt() *FollowUpItemsPrompt {
	return &FollowUpItemsPrompt{}
}
