package reviews

import (
	"github.com/golly-go/golly"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockFeedbackService struct {
	mock.Mock
}

func (m *MockFeedbackService) FindForCode(gctx golly.Context, code string) (Feedback, error) {
	args := m.Called(gctx, code)
	return args.Get(0).(Feedback), args.Error(1)
}

func (m *MockFeedbackService) FindByID(gctx golly.Context, id uuid.UUID) (Feedback, error) {
	args := m.Called(gctx, id)
	return args.Get(0).(Feedback), args.Error(1)
}

func (m *MockFeedbackService) FindByIDs(gctx golly.Context, id uuid.UUIDs) ([]Feedback, error) {
	args := m.Called(gctx, id)
	return args.Get(0).([]Feedback), args.Error(1)
}

func (m *MockFeedbackService) FindByID_Unsafe(gctx golly.Context, id uuid.UUID, scopes ...func(*gorm.DB) *gorm.DB) (Feedback, error) {
	args := m.Called(gctx, id, scopes)
	return args.Get(0).(Feedback), args.Error(1)
}

func (m *MockFeedbackService) FindByIDAndCode_Unsafe(gctx golly.Context, id uuid.UUID, code string) (Feedback, error) {
	args := m.Called(gctx, id, code)
	return args.Get(0).(Feedback), args.Error(1)
}

func (m *MockFeedbackService) FindDetailsByFeedbackID_Unsafe(gctx golly.Context, id uuid.UUID) (FeedbackDetails, error) {
	args := m.Called(gctx, id)
	return args.Get(0).(FeedbackDetails), args.Error(1)
}

func (m *MockFeedbackService) PluckEmailsForSearch(gctx golly.Context, email string) ([]string, error) {
	args := m.Called(gctx, email)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockFeedbackService) FindSummary_Permissioned(gctx golly.Context, feedbackID uuid.UUID) (FeedbackSummary, error) {
	args := m.Called(gctx, feedbackID)
	return args.Get(0).(FeedbackSummary), args.Error(1)
}

var _ ReviewService = &MockFeedbackService{}

func UseMockService(gctx golly.Context, mock *MockFeedbackService) golly.Context {
	gctx.Set(serviceCtxKey, mock)
	return gctx
}
