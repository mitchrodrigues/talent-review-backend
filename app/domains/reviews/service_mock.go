package reviews

import (
	"github.com/golly-go/golly"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockReviewService struct {
	mock.Mock
}

func (m *MockReviewService) FindFeedbackForCode(gctx golly.Context, code string) (Feedback, error) {
	args := m.Called(gctx, code)
	return args.Get(0).(Feedback), args.Error(1)
}

func (m *MockReviewService) FindFeedbackByID(gctx golly.Context, id uuid.UUID) (Feedback, error) {
	args := m.Called(gctx, id)
	return args.Get(0).(Feedback), args.Error(1)
}

func (m *MockReviewService) FindFeedbackByIDAndCode(gctx golly.Context, id uuid.UUID, code string) (Feedback, error) {
	args := m.Called(gctx, id, code)
	return args.Get(0).(Feedback), args.Error(1)
}

func (m *MockReviewService) FindFeedbackDetailsByFeedbackID_Unsafe(gctx golly.Context, id uuid.UUID) (FeedbackDetails, error) {
	args := m.Called(gctx, id)
	return args.Get(0).(FeedbackDetails), args.Error(1)
}

func (m *MockReviewService) FindFeedbackEmailsBySearch(gctx golly.Context, email string) ([]string, error) {
	args := m.Called(gctx, email)
	return args.Get(0).([]string), args.Error(1)
}

var _ ReviewService = &MockReviewService{}

func UseMockService(gctx golly.Context, mock *MockReviewService) golly.Context {
	gctx.Set(serviceCtxKey, mock)
	return gctx
}
