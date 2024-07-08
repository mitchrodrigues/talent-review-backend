package reviews

import (
	"context"
	"testing"
	"time"

	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
	"github.com/golly-go/plugins/orm"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/employees"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/employees/employee"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Define your test function
func TestCreateBulkFeedback(t *testing.T) {
	tests := []struct {
		name           string
		input          CreateBulkFeedbackInput
		mockSetup      func(golly.Context, *employees.MockEmployeeService, *eventsource.MockCommandHandler)
		expectedResult int
		expectedError  error
	}{
		{
			name: "Basic case",
			input: CreateBulkFeedbackInput{
				EmployeeIDs:      []uuid.UUID{uuid.New(), uuid.New()},
				AdditionalEmails: []string{"test1@example.com", "test2@example.com"},
				IncludeTeam:      false,
				CollectionEndAt:  time.Now().Add(24 * time.Hour),
			},
			mockSetup: func(gctx golly.Context, employeeService *employees.MockEmployeeService, callService *eventsource.MockCommandHandler) {
				employeeService.On("FindEmployeesByIDS", mock.Anything, mock.Anything).Return([]employees.Employee{
					{Aggregate: employee.Aggregate{ModelUUID: orm.ModelUUID{ID: uuid.New()}, Email: "employee1@example.com"}},
					{Aggregate: employee.Aggregate{ModelUUID: orm.ModelUUID{ID: uuid.New()}, Email: "employee2@example.com"}},
				}, nil)
				callService.On("Call", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(4)
			},
			expectedResult: 4, // 2 employees x 2 emails each = 4 feedback entries
			expectedError:  nil,
		},
		{
			name: "With duplicate emails",
			input: CreateBulkFeedbackInput{
				EmployeeIDs:      []uuid.UUID{uuid.New(), uuid.New()},
				AdditionalEmails: []string{"employee3@example.com", "employee4@example.com"},
				IncludeTeam:      false,
				CollectionEndAt:  time.Now().Add(24 * time.Hour),
			},
			mockSetup: func(gctx golly.Context, employeeService *employees.MockEmployeeService, callService *eventsource.MockCommandHandler) {
				employeeService.On("FindEmployeesByIDS", mock.Anything, mock.Anything).Return([]employees.Employee{
					{Aggregate: employee.Aggregate{ModelUUID: orm.ModelUUID{ID: uuid.New()}, Email: "employee3@example.com"}},
					{Aggregate: employee.Aggregate{ModelUUID: orm.ModelUUID{ID: uuid.New()}, Email: "employee4@example.com"}},
				}, nil)
				callService.On("Call", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(2)
			},
			expectedResult: 2, // 2 employees x 2 emails each - 1 duplicate each = 2 feedback entries
			expectedError:  nil,
		},
		{
			name: "Include team",
			input: CreateBulkFeedbackInput{
				EmployeeIDs:      []uuid.UUID{uuid.New(), uuid.New()},
				AdditionalEmails: []string{},
				IncludeTeam:      true,
				CollectionEndAt:  time.Now().Add(24 * time.Hour),
			},
			mockSetup: func(gctx golly.Context, employeeService *employees.MockEmployeeService, callService *eventsource.MockCommandHandler) {
				teamID := uuid.New()
				employeeService.On("FindEmployeesByIDS", mock.Anything, mock.Anything).Return([]employees.Employee{
					{Aggregate: employee.Aggregate{ModelUUID: orm.ModelUUID{ID: uuid.New()}, TeamID: &teamID, Email: "employee5@example.com"}},
					{Aggregate: employee.Aggregate{ModelUUID: orm.ModelUUID{ID: uuid.New()}, TeamID: &teamID, Email: "employee6@example.com"}},
				}, nil)
				employeeService.On("FindEmployeesForTeam", mock.Anything, mock.Anything, mock.Anything).Return([]employees.Employee{
					{Aggregate: employee.Aggregate{ModelUUID: orm.ModelUUID{ID: uuid.New()}, Email: "teammate1@example.com"}},
					{Aggregate: employee.Aggregate{ModelUUID: orm.ModelUUID{ID: uuid.New()}, Email: "teammate2@example.com"}},
				}, nil)
				callService.On("Call", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(4)
			},
			expectedResult: 4, // 2 employees x (2 teammates) = 4 feedback entries
			expectedError:  nil,
		},
		{
			name: "Include All Together",
			input: CreateBulkFeedbackInput{
				EmployeeIDs:      []uuid.UUID{uuid.New(), uuid.New()},
				AdditionalEmails: []string{"employee3@example.com", "employee4@example.com", "test1@example.com", "test2@example.com", "teammate1@example.com"},
				IncludeTeam:      true,
				CollectionEndAt:  time.Now().Add(24 * time.Hour),
			},
			mockSetup: func(gctx golly.Context, employeeService *employees.MockEmployeeService, callService *eventsource.MockCommandHandler) {
				teamID := uuid.New()
				employeeService.On("FindEmployeesByIDS", mock.Anything, mock.Anything).Return([]employees.Employee{
					{Aggregate: employee.Aggregate{ModelUUID: orm.ModelUUID{ID: uuid.New()}, TeamID: &teamID, Email: "employee5@example.com"}},
					{Aggregate: employee.Aggregate{ModelUUID: orm.ModelUUID{ID: uuid.New()}, TeamID: &teamID, Email: "employee6@example.com"}},
				}, nil)
				employeeService.On("FindEmployeesForTeam", mock.Anything, mock.Anything, mock.Anything).Return([]employees.Employee{
					{Aggregate: employee.Aggregate{ModelUUID: orm.ModelUUID{ID: uuid.New()}, Email: "teammate1@example.com"}},
					{Aggregate: employee.Aggregate{ModelUUID: orm.ModelUUID{ID: uuid.New()}, Email: "teammate2@example.com"}},
				}, nil)
				callService.On("Call", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(12)
			},
			expectedResult: 12, // 2 employees x (2 teammates) = 4 feedback entries
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockEmployeeService := new(employees.MockEmployeeService)
			mockHandler := new(eventsource.MockCommandHandler)

			gctx := golly.NewContext(context.TODO())
			gctx = orm.CreateTestContext(gctx)

			employees.UseMockService(gctx, mockEmployeeService)
			eventsource.UseMockHandler(gctx, mockHandler)

			// Set up mocks
			tt.mockSetup(gctx, mockEmployeeService, mockHandler)

			// Call the function under test
			result, err := CreateBulkFeedback(gctx, tt.input, eventsource.Metadata{})

			// Assertions
			assert.Equal(t, tt.expectedError, err)
			assert.Len(t, result, tt.expectedResult)

			// Assert that the expectations were met
			mockEmployeeService.AssertExpectations(t)
			mockHandler.AssertExpectations(t)
		})
	}
}
