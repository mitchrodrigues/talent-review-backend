package reviews

import (
	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/orm"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/common"
)

const (
	serviceCtxKey golly.ContextKeyT = "feedbackService"
)

type ReviewService interface {
	FindFeedbackForCode(gctx golly.Context, code string) (Feedback, error)
	FindFeedbackByID(gctx golly.Context, id uuid.UUID) (Feedback, error)
	FindFeedbackByIDAndCode(gctx golly.Context, id uuid.UUID, code string) (Feedback, error)
	FindFeedbackDetailsByFeedbackID_Unsafe(gctx golly.Context, id uuid.UUID) (FeedbackDetails, error)
	FindFeedbackEmailsBySearch(gctx golly.Context, email string) ([]string, error)
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

func (DefaultReviewService) FindFeedbackByID(gctx golly.Context, id uuid.UUID) (Feedback, error) {
	var feedback Feedback

	err := orm.
		DB(gctx).
		Model(feedback).
		Find(&feedback, "id = ?", id).
		Error

	return feedback, err
}

func (DefaultReviewService) FindFeedbackByIDAndCode(gctx golly.Context, id uuid.UUID, code string) (Feedback, error) {
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
