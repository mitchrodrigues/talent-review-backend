package reviews

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/employees"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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
				managerID := uuid.New()
				organizationID := uuid.New()
				employeeService.On("FindEmployeeByUserID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(employees.NewTestEmployee(managerID, organizationID, "manager@example.com", nil), nil)
				employeeService.On("FindEmployeesByManagerAndIDS", mock.Anything, managerID, mock.Anything).Return([]employees.Employee{
					employees.NewTestEmployee(uuid.New(), organizationID, "employee1@example.com", nil),
					employees.NewTestEmployee(uuid.New(), organizationID, "employee2@example.com", nil),
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
				managerID := uuid.New()
				organizationID := uuid.New()
				employeeService.On("FindEmployeeByUserID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(employees.NewTestEmployee(managerID, organizationID, "manager@example.com", nil), nil)
				employeeService.On("FindEmployeesByManagerAndIDS", mock.Anything, managerID, mock.Anything).Return([]employees.Employee{
					employees.NewTestEmployee(uuid.New(), organizationID, "employee3@example.com", nil),
					employees.NewTestEmployee(uuid.New(), organizationID, "employee4@example.com", nil),
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
				managerID := uuid.New()
				organizationID := uuid.New()
				teamID := uuid.New()
				employeeService.On("FindEmployeeByUserID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(employees.NewTestEmployee(managerID, organizationID, "manager@example.com", nil), nil)
				employeeService.On("FindEmployeesByManagerAndIDS", mock.Anything, managerID, mock.Anything).Return([]employees.Employee{
					employees.NewTestEmployeeWithTeam(uuid.New(), organizationID, "employee5@example.com", teamID),
					employees.NewTestEmployeeWithTeam(uuid.New(), organizationID, "employee6@example.com", teamID),
				}, nil)
				employeeService.On("FindEmployeesForTeam", mock.Anything, teamID, mock.Anything).Return([]employees.Employee{
					employees.NewTestEmployee(uuid.New(), organizationID, "teammate1@example.com", nil),
					employees.NewTestEmployee(uuid.New(), organizationID, "teammate2@example.com", nil),
				}, nil)
				callService.On("Call", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(4)
			},
			expectedResult: 4, // 2 employees x 2 teammates = 4 feedback entries
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
				managerID := uuid.New()
				organizationID := uuid.New()
				teamID := uuid.New()
				employeeService.On("FindEmployeeByUserID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(employees.NewTestEmployee(managerID, organizationID, "manager@example.com", nil), nil)
				employeeService.On("FindEmployeesByManagerAndIDS", mock.Anything, managerID, mock.Anything).Return([]employees.Employee{
					employees.NewTestEmployeeWithTeam(uuid.New(), organizationID, "employee5@example.com", teamID),
					employees.NewTestEmployeeWithTeam(uuid.New(), organizationID, "employee6@example.com", teamID),
				}, nil)
				employeeService.On("FindEmployeesForTeam", mock.Anything, teamID, mock.Anything).Return([]employees.Employee{
					employees.NewTestEmployee(uuid.New(), organizationID, "teammate1@example.com", nil),
					employees.NewTestEmployee(uuid.New(), organizationID, "teammate2@example.com", nil),
				}, nil)
				callService.On("Call", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(12)
			},
			expectedResult: 12, // Multiple combinations of employees and emails including team and additional emails
			expectedError:  nil,
		},
		{
			name: "Manager not found",
			input: CreateBulkFeedbackInput{
				EmployeeIDs:      []uuid.UUID{uuid.New(), uuid.New()},
				AdditionalEmails: []string{"test1@example.com", "test2@example.com"},
				IncludeTeam:      false,
				CollectionEndAt:  time.Now().Add(24 * time.Hour),
			},
			mockSetup: func(gctx golly.Context, employeeService *employees.MockEmployeeService, callService *eventsource.MockCommandHandler) {
				employeeService.On("FindEmployeeByUserID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(employees.Employee{}, fmt.Errorf("manager not found"))
			},
			expectedResult: 0,
			expectedError:  fmt.Errorf("you are not a manager of any team"),
		},
		{
			name: "No employees found",
			input: CreateBulkFeedbackInput{
				EmployeeIDs:      []uuid.UUID{uuid.New(), uuid.New()},
				AdditionalEmails: []string{"test1@example.com", "test2@example.com"},
				IncludeTeam:      false,
				CollectionEndAt:  time.Now().Add(24 * time.Hour),
			},
			mockSetup: func(gctx golly.Context, employeeService *employees.MockEmployeeService, callService *eventsource.MockCommandHandler) {
				managerID := uuid.New()
				organizationID := uuid.New()
				employeeService.On("FindEmployeeByUserID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(employees.NewTestEmployee(managerID, organizationID, "manager@example.com", nil), nil)
				employeeService.On("FindEmployeesByManagerAndIDS", mock.Anything, managerID, mock.Anything).Return([]employees.Employee{}, nil)
			},
			expectedResult: 0,
			expectedError:  fmt.Errorf("you do not have any employees matching the criteria"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockEmployeeService := new(employees.MockEmployeeService)
			mockHandler := new(eventsource.MockCommandHandler)

			gctx := golly.NewContext(context.TODO())

			employees.UseMockService(gctx, mockEmployeeService)
			eventsource.UseMockHandler(gctx, mockHandler)

			// Set up mocks
			tt.mockSetup(gctx, mockEmployeeService, mockHandler)

			// Call the function under test
			result, err := CreateBulkFeedback(gctx, tt.input, eventsource.Metadata{})

			// Assertions
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Len(t, result, tt.expectedResult)

			// Assert that the expectations were met
			mockEmployeeService.AssertExpectations(t)
			mockHandler.AssertExpectations(t)
		})
	}
}
