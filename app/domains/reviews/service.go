package reviews

import (
	"fmt"

	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/orm"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/common"
	"gorm.io/gorm"
)

const (
	serviceCtxKey golly.ContextKeyT = "feedbackService"
)

type ReviewService interface {
	FindFeedbackForCode(gctx golly.Context, code string) (Feedback, error)
	FindFeedbackByID(gctx golly.Context, id uuid.UUID, scopes ...func(*gorm.DB) *gorm.DB) (Feedback, error)

	FindFeedbackEmailsBySearch(gctx golly.Context, email string) ([]string, error)
	FindFeedbackSummary_Permissioned(gctx golly.Context, feedbackID uuid.UUID) (FeedbackSummary, error)

	FindFeedbackByID_Unsafe(gctx golly.Context, id uuid.UUID) (Feedback, error)
	FindFeedbackByIDAndCode_Unsafe(gctx golly.Context, id uuid.UUID, code string) (Feedback, error)
	FindFeedbackDetailsByFeedbackID_Unsafe(gctx golly.Context, id uuid.UUID) (FeedbackDetails, error)
}

type DefaultReviewService struct{}

func (DefaultReviewService) FindFeedbackForCode(gctx golly.Context, code string) (Feedback, error) {
	var feedback Feedback

	err := orm.
		DB(gctx).
		Model(feedback).
		Find(&feedback, "code = ?", code).
		Error

	return feedback, err
}

func (DefaultReviewService) FindFeedbackByID(gctx golly.Context, id uuid.UUID, scopes ...func(*gorm.DB) *gorm.DB) (Feedback, error) {
	var feedback Feedback

	err := orm.
		DB(gctx).
		Model(feedback).
		Scopes(common.OrganizationIDScopeForContext(gctx, "feedbacks")).
		Scopes(scopes...).
		Find(&feedback, "feedbacks.id = ?", id).
		Error

	return feedback, err
}

func (DefaultReviewService) FindFeedbackByID_Unsafe(gctx golly.Context, id uuid.UUID) (Feedback, error) {
	var feedback Feedback

	err := orm.
		DB(gctx).
		Model(feedback).
		Find(&feedback, "feedbacks.id = ?", id).
		Error

	return feedback, err
}

func (DefaultReviewService) FindFeedbackSummary_Permissioned(gctx golly.Context, feedbackID uuid.UUID) (FeedbackSummary, error) {
	return golly.LoadData(
		gctx,
		fmt.Sprintf("feedbackSummaries:%s", feedbackID),
		func(gctx golly.Context) (FeedbackSummary, error) {
			var summary FeedbackSummary

			err := orm.
				DB(gctx).
				Model(summary).
				Scopes(
					common.OrganizationIDScopeForContext(gctx, "feedback_summaries"),
					common.UserIsManagerScope(gctx, "feedback_summaries"),
				).
				Find(&summary, "feedback_id = ?", feedbackID).
				Error

			return summary, err
		})
}

func (DefaultReviewService) FindFeedbackByIDAndCode_Unsafe(gctx golly.Context, id uuid.UUID, code string) (Feedback, error) {
	var feedback Feedback

	err := orm.
		DB(gctx).
		Model(feedback).
		Find(&feedback, "id = ? AND code = ?", id, code).
		Error

	return feedback, err
}

func (DefaultReviewService) FindFeedbackDetailsByFeedbackID_Unsafe(gctx golly.Context, id uuid.UUID) (FeedbackDetails, error) {
	var details FeedbackDetails

	err := orm.
		DB(gctx).
		Model(details).
		Find(&details, "feedback_id = ?", id).
		Error

	return details, err
}

func (DefaultReviewService) FindFeedbackEmailsBySearch(gctx golly.Context, email string) ([]string, error) {
	var emails []string

	err := orm.
		DB(gctx).
		Model(&Feedback{}).
		Scopes(common.OrganizationIDScopeForContext(gctx)).
		Select("DISTINCT(email) AS email").
		Where("LOWER(email) LIKE ?", email+"%").
		Pluck("email", &emails).
		Error

	return emails, err
}

func Service(gctx golly.Context) ReviewService {
	if service, ok := gctx.Get(serviceCtxKey); ok {
		return service.(ReviewService)
	}

	service := DefaultReviewService{}
	gctx.Set(serviceCtxKey, service)

	return service
}

var _ ReviewService = DefaultReviewService{}
