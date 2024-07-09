package reviews

import (
	"context"
	"testing"

	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/orm"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/reviews/feedback"
	"github.com/stretchr/testify/assert"
)

// Test for FindFeedbackForCode
func TestFindFeedbackForCode(t *testing.T) {
	gctx := orm.CreateTestContext(golly.NewContext(context.TODO()), Feedback{}, FeedbackDetails{})

	// Seed the database
	fb := Feedback{
		Aggregate: feedback.Aggregate{
			ModelUUID: orm.NewModelUUID(),
			Code:      "test-code",
			Email:     "test@example.com",
		},
	}

	orm.DB(gctx).Create(&fb)

	service := DefaultReviewService{}
	result, err := service.FindFeedbackForCode(gctx, "test-code")

	assert.NoError(t, err)
	assert.Equal(t, fb.Code, result.Code)
	assert.Equal(t, fb.Email, result.Email)
}

// Test for FindFeedbackByID
func TestFindFeedbackByID(t *testing.T) {
	gctx := orm.CreateTestContext(golly.NewContext(context.TODO()), Feedback{}, FeedbackDetails{})

	// Seed the database
	fb := Feedback{
		Aggregate: feedback.Aggregate{
			ModelUUID: orm.NewModelUUID(),
			Code:      "test-code",
			Email:     "test@example.com",
		},
	}

	orm.DB(gctx).Create(&fb)

	service := DefaultReviewService{}
	result, err := service.FindFeedbackByID(gctx, fb.ID)

	assert.NoError(t, err)
	assert.Equal(t, fb.ID, result.ID)
	assert.Equal(t, fb.Email, result.Email)
}

// Test for FindFeedbackByIDAndCode
func TestFindFeedbackByIDAndCode(t *testing.T) {
	gctx := orm.CreateTestContext(golly.NewContext(context.TODO()), Feedback{}, FeedbackDetails{})

	// Seed the database
	fb := Feedback{
		Aggregate: feedback.Aggregate{
			ModelUUID: orm.NewModelUUID(),
			Code:      "test-code",
			Email:     "test@example.com",
		},
	}

	orm.DB(gctx).Create(&fb)

	service := DefaultReviewService{}
	result, err := service.FindFeedbackByIDAndCode_Unsafe(gctx, fb.ID, "test-code")

	assert.NoError(t, err)
	assert.Equal(t, fb.ID, result.ID)
	assert.Equal(t, fb.Code, result.Code)
}

// Test for FindFeedbackDetailsByFeedbackID_Unsafe
func TestFindFeedbackDetailsByFeedbackID_Unsafe(t *testing.T) {
	gctx := orm.CreateTestContext(golly.NewContext(context.TODO()), Feedback{}, FeedbackDetails{})

	// Seed the database
	fb := Feedback{
		Aggregate: feedback.Aggregate{
			ModelUUID: orm.NewModelUUID(),
			Code:      "test-code",
			Email:     "test@example.com",
		},
	}
	orm.DB(gctx).Create(&fb)

	details := FeedbackDetails{
		FeedbackDetails: feedback.FeedbackDetails{
			ModelUUID:  orm.NewModelUUID(),
			FeedbackID: fb.ID,
		},
	}

	orm.DB(gctx).Create(&details)

	service := DefaultReviewService{}
	result, err := service.FindFeedbackDetailsByFeedbackID_Unsafe(gctx, fb.ID)

	assert.NoError(t, err)
	assert.Equal(t, details.FeedbackID, result.FeedbackID)
}

// Test for FindFeedbackEmailsBySearch
func TestFindFeedbackEmailsBySearch(t *testing.T) {
	gctx := orm.CreateTestContext(golly.NewContext(context.TODO()), Feedback{}, FeedbackDetails{})

	// Seed the database
	feedback1 := Feedback{
		Aggregate: feedback.Aggregate{
			ModelUUID: orm.NewModelUUID(),
			Code:      "code1",
			Email:     "test1@example.com",
		},
	}
	feedback2 := Feedback{

		Aggregate: feedback.Aggregate{
			ModelUUID: orm.NewModelUUID(),
			Code:      "code2",
			Email:     "test2@example.com",
		},
	}
	orm.DB(gctx).Create(&feedback1)
	orm.DB(gctx).Create(&feedback2)

	service := DefaultReviewService{}
	result, err := service.FindFeedbackEmailsBySearch(gctx, "test")

	assert.NoError(t, err)
	assert.Contains(t, result, "test1@example.com")
	assert.Contains(t, result, "test2@example.com")
}
