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
	FindForCode(gctx golly.Context, code string) (Feedback, error)
	FindByID(gctx golly.Context, id uuid.UUID) (Feedback, error)
	FindByIDs(gctx golly.Context, id uuid.UUIDs) ([]Feedback, error)

	PluckEmailsForSearch(gctx golly.Context, email string) ([]string, error)
	FindSummary_Permissioned(gctx golly.Context, feedbackID uuid.UUID) (FeedbackSummary, error)
	FindAll_Permissioned(gctx golly.Context, scopes ...func(*gorm.DB) *gorm.DB) ([]Feedback, error)

	FindByID_Unsafe(gctx golly.Context, id uuid.UUID, scopes ...func(*gorm.DB) *gorm.DB) (Feedback, error)
	FindByIDAndCode_Unsafe(gctx golly.Context, id uuid.UUID, code string) (Feedback, error)

	FindDetailsByFeedbackID_Unsafe(gctx golly.Context, id uuid.UUID) (FeedbackDetails, error)
}

type DefaultReviewService struct{}

func (DefaultReviewService) FindForCode(gctx golly.Context, code string) (Feedback, error) {
	var feedback Feedback

	err := orm.
		DB(gctx).
		Model(feedback).
		Find(&feedback, "code = ?", code).
		Error

	return feedback, err
}

func (service DefaultReviewService) FindByID(gctx golly.Context, id uuid.UUID) (Feedback, error) {
	return service.FindByID_Unsafe(gctx, id, common.OrganizationIDScopeForContext(gctx, "feedbacks"))
}

func (service DefaultReviewService) FindAll_Permissioned(gctx golly.Context, scopes ...func(*gorm.DB) *gorm.DB) ([]Feedback, error) {
	var feedbacks []Feedback

	err := orm.
		DB(gctx).
		Scopes(
			common.OrganizationIDScopeForContext(gctx, "feedbacks"),
			common.JoinUserEmployeeRecord(gctx)).
		Scopes(scopes...).
		Joins("JOIN employees employee ON employee.id = feedbacks.employee_id").
		Where("user_employee_record.id = employee.manager_id OR feedbacks.email = user_employee_record.email").
		Find(&feedbacks).
		Error

	return feedbacks, err
}

func (DefaultReviewService) FindByIDs(gctx golly.Context, ids uuid.UUIDs) ([]Feedback, error) {
	var feedbacks []Feedback

	err := orm.
		DB(gctx).
		Scopes(common.OrganizationIDScopeForContext(gctx, "feedbacks")).
		Find(&feedbacks, "id IN ?", ids).
		Error

	return feedbacks, err
}

func (DefaultReviewService) FindSummary_Permissioned(gctx golly.Context, feedbackID uuid.UUID) (FeedbackSummary, error) {
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

func (DefaultReviewService) PluckEmailsForSearch(gctx golly.Context, email string) ([]string, error) {
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

func (DefaultReviewService) FindDetailsByFeedbackID_Unsafe(gctx golly.Context, id uuid.UUID) (FeedbackDetails, error) {
	var details FeedbackDetails

	err := orm.
		DB(gctx).
		Model(details).
		Find(&details, "feedback_id = ?", id).
		Error

	return details, err
}

func (DefaultReviewService) FindByID_Unsafe(gctx golly.Context, id uuid.UUID, scopes ...func(*gorm.DB) *gorm.DB) (Feedback, error) {
	var feedback Feedback

	err := orm.
		DB(gctx).
		Model(feedback).
		Scopes(scopes...).
		Find(&feedback, "feedbacks.id = ?", id).
		Error

	return feedback, err
}

func (service DefaultReviewService) FindByIDAndCode_Unsafe(gctx golly.Context, id uuid.UUID, code string) (Feedback, error) {
	return service.FindByID_Unsafe(gctx, id, func(db *gorm.DB) *gorm.DB {
		return db.Where("code = ?", code)
	})
}

func FeedbackService(gctx golly.Context) ReviewService {
	if service, ok := gctx.Get(serviceCtxKey); ok {
		return service.(ReviewService)
	}

	service := DefaultReviewService{}
	gctx.Set(serviceCtxKey, service)

	return service
}

var _ ReviewService = DefaultReviewService{}
