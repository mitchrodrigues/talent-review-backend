package reviews

import "github.com/mitchrodrigues/talent-review-backend/app/domains/reviews/feedback"

type Feedback struct {
	feedback.Aggregate
}

func (Feedback) TableName() string { return "feedbacks" }

type FeedbackDetails struct {
	feedback.FeedbackDetails
}

func (FeedbackDetails) TableName() string { return "feedback_details" }
