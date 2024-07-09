package tara

import (
	"github.com/golly-go/golly"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/openai"
)

func Generate[T openai.CompletionPrompt](gctx golly.Context, prompt T, aiContexts ...openai.AIContext) error {
	_, err := openai.Completion(gctx, openai.LLM(gctx), prompt)
	return err
}

func Initailizer(app golly.Application) error {
	return openai.Initailizer(openai.Config{
		Token: app.Config.GetString("openai.token"),
		AIPretendsToBe: "Pretend you are a performance management assistant providing insights and " +
			"recommendations for team growth, performance reviews, and inclusivity.",
		AIScenarioContext: "You are assisting in the Talent Radar application, focusing on team growth " +
			"and performance metrics. Your tasks include generating insights from performance data, offering " +
			"personalized development suggestions, reviewing feedback for inclusivity, predicting future " +
			"performance trends, and ensuring fair and unbiased performance reviews.",
	})
}
