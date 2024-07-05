package reviews

import (
	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/orm"
	"github.com/google/uuid"
)

func FindFeedbackForCode(gctx golly.Context, code string) (Feedback, error) {
	var feedback Feedback

	err := orm.
		DB(gctx).
		Model(feedback).
		Find(&feedback, "code = ?", code).
		Error

	return feedback, err
}

func FindFeedbackByID(gctx golly.Context, id uuid.UUID) (Feedback, error) {
	var feedback Feedback

	err := orm.
		DB(gctx).
		Model(feedback).
		Find(&feedback, "id = ?", id).
		Error

	return feedback, err
}

func FindFeedbackByIDAndCode(gctx golly.Context, id uuid.UUID, code string) (Feedback, error) {
	var feedback Feedback

	err := orm.
		DB(gctx).
		Model(feedback).
		Find(&feedback, "id = ? AND code = ?", id, code).
		Error

	return feedback, err
}

func FindFeedbackDetailsByFeedbackID_Unsafe(gctx golly.Context, id uuid.UUID) (FeedbackDetails, error) {
	var details FeedbackDetails

	err := orm.
		DB(gctx).
		Model(details).
		Find(&details, "feedback_id = ?", id).
		Error

	return details, err
}
