package reviews

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/gql"
	"github.com/golly-go/plugins/orm"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/employees"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/reviews/feedback"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/esbackend"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/identity"
	"github.com/stretchr/testify/assert"
)

// Create a test context
func createTestContext() golly.Context {
	return orm.CreateTestContext(golly.NewContext(context.TODO()),
		Feedback{},
		FeedbackDetails{},
		esbackend.Event{},
		employees.Employee{},
		employees.Team{})
}

func formatTimestampGQL(t time.Time) string {
	return t.Format(time.RFC3339Nano)
}

func TestFeedbacks(t *testing.T) {
	type FeedbackResponse struct {
		ID              uuid.UUID `json:"id"`
		CreatedAt       string    `json:"createdAt"`
		CollectionEndAt string    `json:"collectionEndAt"`
		Email           string    `json:"email"`
		SubmittedAt     *string   `json:"submittedAt"`
		Employee        struct {
			ID   uuid.UUID `json:"id"`
			Name string    `json:"name"`
		} `json:"employee"`
	}

	type GroupedFeedbackResponse struct {
		Feedbacks []FeedbackResponse `json:"feedbacks"`
	}

	type testCase struct {
		name      string
		identity  identity.Identity
		variables map[string]interface{}
		expected  GroupedFeedbackResponse
		hasError  bool
	}

	// Set up the test context and seed the database
	gctx := createTestContext()

	// Seed the database with feedback entries
	managerUserID := uuid.New()
	organizationID := uuid.New()

	manager := employees.NewTestEmployee(uuid.New(), organizationID, "manager@example.com", &managerUserID)
	orm.DB(gctx).Create(&manager)

	team := employees.NewTestTeam(uuid.New(), manager.OrganizationID, manager.ID)
	orm.DB(gctx).Create(&team)

	employee := employees.NewTestEmployeeWithTeam(uuid.New(), organizationID, "testemployee@example.com", team.ID)

	orm.DB(gctx).Create(&employee)

	fb1 := Feedback{
		Aggregate: feedback.Aggregate{
			ModelUUID: orm.ModelUUID{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
			},
			Code:            "test-code",
			Email:           "test@example.com",
			CollectionEndAt: time.Now().Add(7 * 24 * time.Hour),
			EmployeeID:      employee.ID,
			OrganizationID:  organizationID,
		},
	}
	orm.DB(gctx).Create(&fb1)

	fb2 := Feedback{
		Aggregate: feedback.Aggregate{
			ModelUUID: orm.ModelUUID{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
			},
			Code:            "test-code-2",
			Email:           "test@example.com",
			CollectionEndAt: time.Now().Add(7 * 24 * time.Hour),
			EmployeeID:      employee.ID,
			OrganizationID:  organizationID,
		},
	}
	orm.DB(gctx).Create(&fb2)

	// Define test cases
	testCases := []testCase{
		{
			name: "Valid feedbacks query, logged in",
			identity: identity.Identity{
				UID:            managerUserID,
				OrganizationID: organizationID,
			},
			variables: map[string]interface{}{
				"pagination": map[string]interface{}{
					"first": 10,
				},
			},
			expected: GroupedFeedbackResponse{
				Feedbacks: []FeedbackResponse{
					{
						ID:              fb1.ID,
						CreatedAt:       formatTimestampGQL(fb1.CreatedAt),
						CollectionEndAt: formatTimestampGQL(fb1.CollectionEndAt),
						Email:           fb1.Email,
						SubmittedAt:     nil,
						Employee: struct {
							ID   uuid.UUID `json:"id"`
							Name string    `json:"name"`
						}{
							ID:   employee.ID,
							Name: employee.Name,
						},
					},

					{
						ID:              fb2.ID,
						CreatedAt:       formatTimestampGQL(fb2.CreatedAt),
						CollectionEndAt: formatTimestampGQL(fb2.CollectionEndAt),
						Email:           fb2.Email,
						SubmittedAt:     nil,
						Employee: struct {
							ID   uuid.UUID `json:"id"`
							Name string    `json:"name"`
						}{
							ID:   employee.ID,
							Name: employee.Name,
						},
					},
				},
			},
		},
		{
			name: "User not logged in",
			identity: identity.Identity{
				UID:            uuid.Nil,
				OrganizationID: uuid.Nil,
			},
			variables: map[string]interface{}{
				"pagination": map[string]interface{}{
					"first": 10,
				},
			},
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set the identity in the context
			ctx := identity.ToContext(gctx, tc.identity)

			// Define the test query with variables
			query := `
				query feedbacks($pagination: PaginationInput) {
					groupedFeedbacks(pagination: $pagination) {
						feedbacks {
							id
							createdAt
							collectionEndAt
							email
							submittedAt
							employee {
								id
								name
							}
						}
					}
				}
			`

			r, err := gql.ExecuteGraphQLQuery(ctx, queries, query, tc.variables)
			assert.NotNil(t, r)

			if tc.hasError {
				assert.NotNil(t, r.Errors)
				return
			}

			assert.NoError(t, err)
			assert.Nil(t, r.Errors)
			assert.NotNil(t, r.Data)

			fmt.Printf("%#v\n", r)

			b, _ := json.Marshal(r.Data.(map[string]interface{})["groupedFeedbacks"])

			var gresults []GroupedFeedbackResponse

			err = json.Unmarshal(b, &gresults)
			assert.NoError(t, err)

			assert.GreaterOrEqual(t, len(gresults), 1)

			for _, results := range gresults {
				assert.Equal(t, len(tc.expected.Feedbacks), len(results.Feedbacks))

				for _, expected := range tc.expected.Feedbacks {
					result := golly.Find(results.Feedbacks, func(edge FeedbackResponse) bool {
						return edge.ID == expected.ID
					})

					assert.NotNil(t, result)

					assert.Equal(t, expected.ID, result.ID)
					assert.Equal(t, expected.CreatedAt, result.CreatedAt)
					assert.Equal(t, expected.CollectionEndAt, result.CollectionEndAt)
					assert.Equal(t, expected.Email, result.Email)
					assert.Equal(t, expected.SubmittedAt, result.SubmittedAt)
					assert.Equal(t, expected.Employee.ID, result.Employee.ID)
					assert.Equal(t, expected.Employee.Name, result.Employee.Name)
				}

			}
		})
	}
}

func TestFeedbackForCode(t *testing.T) {

	timestamp := time.Now()

	type FeedbackResponse struct {
		ID uuid.UUID `json:"id"`
	}

	type testCase struct {
		name     string
		code     string
		identity identity.Identity
		expected FeedbackResponse
		hasError bool
	}

	// Set up the test context and seed the database
	gctx := createTestContext()

	// Seed the database with a feedback entry
	fb := Feedback{
		Aggregate: feedback.Aggregate{
			ModelUUID: orm.ModelUUID{
				ID:        uuid.New(),
				CreatedAt: timestamp,
			},
			Code:            "test-code",
			Email:           "test@example.com",
			CollectionEndAt: timestamp.Add(7 * 24 * time.Hour),
			EmployeeID:      uuid.New(),
		},
	}
	orm.DB(gctx).Create(&fb)

	// Seed the database with an employee entry
	employee := employees.NewTestEmployee(fb.EmployeeID, uuid.Nil, "testemployee@examlpe.com", nil)
	orm.DB(gctx).Create(&employee)

	// Define test cases
	testCases := []testCase{
		{
			name: "Valid feedback query, logged in",
			code: fb.Code,
			identity: identity.Identity{
				UID:            uuid.New(),
				OrganizationID: uuid.New(),
			},
			expected: FeedbackResponse{
				ID: fb.ID,
			},
			hasError: false,
		},
		{
			name: "Valid feedback query, logged out",
			code: fb.Code,
			expected: FeedbackResponse{
				ID: fb.ID,
			},
			hasError: false,
		},
		{
			name:     "Invalid feedback code",
			code:     "invalid-code",
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := identity.ToContext(gctx, tc.identity)

			// Define the test query with variables
			query := `
				query feedbackForCode($code: String!) {
					feedbackForCode(code: $code) {
						id
					}
				}
			`

			variables := map[string]interface{}{
				"code": tc.code,
			}

			r, err := gql.ExecuteGraphQLQuery(ctx, queries, query, variables)
			assert.NotNil(t, r)

			if tc.hasError {
				assert.NotNil(t, r.Errors)
				return
			}

			assert.NoError(t, err)
			assert.Nil(t, r.Errors)

			b, _ := json.Marshal(r.Data.(map[string]interface{})["feedbackForCode"])

			var results FeedbackResponse

			err = json.Unmarshal(b, &results)
			assert.NoError(t, err)

			assert.Equal(t, tc.expected.ID, results.ID)
		})
	}
}

func TestUpdateFeedbackDetails_Integration(t *testing.T) {
	type UpdateDetails struct {
		ID    uuid.UUID `json:"id"`
		Email string    `json:"email"`
	}

	type testCase struct {
		name      string
		variables map[string]interface{}
		expected  UpdateDetails
		hasError  bool
		identity  identity.Identity
	}

	// Set up the test context and seed the database
	gctx := createTestContext()

	fb := Feedback{
		Aggregate: feedback.Aggregate{
			ModelUUID:  orm.NewModelUUID(),
			Code:       "test-code",
			Email:      "test@example.com",
			EmployeeID: uuid.New(),
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

	// Define test cases
	testCases := []testCase{
		{
			name:     "Valid update Logged In",
			identity: identity.Identity{UID: uuid.New()},
			variables: map[string]interface{}{
				"code": fb.Code,
				"id":   fb.ID.String(),
				"input": map[string]interface{}{
					"strengths":     "strength",
					"opportunities": "opportunity",
					"additional":    "additional",
					"rating":        5,
					"enoughData":    true,
				},
			},
			expected: UpdateDetails{
				ID:    fb.ID,
				Email: fb.Email,
			},
			hasError: false,
		},
		{
			name: "Valid update Logged Out",
			variables: map[string]interface{}{
				"code": fb.Code,
				"id":   fb.ID.String(),
				"input": map[string]interface{}{
					"strengths":     "strength",
					"opportunities": "opportunity",
					"additional":    "additional",
					"rating":        5,
					"enoughData":    true,
				},
			},
			expected: UpdateDetails{
				ID:    fb.ID,
				Email: fb.Email,
			},
			hasError: false,
		},
		{
			name: "Invalid ID",
			variables: map[string]interface{}{
				"id": uuid.New().String(),
				"input": map[string]interface{}{
					"strengths":     "strength",
					"opportunities": "opportunity",
					"additional":    "additional",
					"rating":        5,
					"enoughData":    true,
				},
			},
			hasError: true,
		},
		{
			name: "Invalid code",
			variables: map[string]interface{}{
				"code": "invalid-code",
				"id":   fb.ID.String(),
				"input": map[string]interface{}{
					"strengths":     "strength",
					"opportunities": "opportunity",
					"additional":    "additional",
					"rating":        5,
					"enoughData":    true,
				},
			},
			hasError: true,
		},
		{
			name: "Missing required input field",
			variables: map[string]interface{}{
				"code": fb.Code,
				"id":   fb.ID.String(),
			},
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Define the test mutation with variables
			mutation := `
				mutation updateFeedbackDetails($id: String!, $code: String!, $input: UpdateFeedbackDetailsInput!) {
					updateFeedbackDetails(id: $id, code: $code, input: $input) {
						id
						email
					}
				}
			`

			gctx := identity.ToContext(gctx, tc.identity)

			r, err := gql.ExecuteGraphQLMutation(gctx, mutations, mutation, tc.variables)
			assert.NotNil(t, r)

			if tc.hasError {
				assert.NotNil(t, r.Errors)
				return
			}

			assert.NoError(t, err)
			assert.Nil(t, r.Errors)

			b, _ := json.Marshal(r.Data.(map[string]interface{})["updateFeedbackDetails"])

			var results UpdateDetails

			err = json.Unmarshal(b, &results)
			assert.NoError(t, err)

			assert.Equal(t, tc.expected.ID, results.ID)
			assert.Equal(t, tc.expected.Email, results.Email)
		})
	}
}

func TestCreateFeedbacks_Integration(t *testing.T) {
	type CreateFeedback struct {
		Email string `json:"email"`
	}

	type testCase struct {
		name      string
		identity  identity.Identity
		variables map[string]interface{}
		expected  []CreateFeedback
		hasError  bool
		empty     bool
		errorMsg  string
	}

	// Set up the test context and seed the database
	gctx := createTestContext()

	// User ID for the manager
	managerUserID := uuid.New()
	organizationID := uuid.New()

	manager := employees.NewTestEmployee(uuid.New(), organizationID, "manager@example.com", &managerUserID)
	orm.DB(gctx).Create(&manager)

	// Team with the manager
	team := employees.NewTestTeam(uuid.New(), organizationID, manager.ID)
	orm.DB(gctx).Create(&team)

	// Employees in the team
	emps := make([]employees.Employee, 2)
	for pos := range emps {
		emps[pos] = employees.NewTestEmployeeWithTeam(uuid.New(), organizationID, fmt.Sprintf("testemployee+%d@example.com", pos), team.ID)
		orm.DB(gctx).Create(&emps[pos])
	}

	// Define test cases
	testCases := []testCase{
		{
			name: "Manager creating feedbacks for team members",
			identity: identity.Identity{
				UID:            managerUserID,
				OrganizationID: organizationID,
			},
			variables: map[string]interface{}{
				"input": map[string]interface{}{
					"employeeIDs":     []string{emps[0].ID.String(), emps[1].ID.String()},
					"includeTeam":     true,
					"collectionEndAt": time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339),
					"additionalEmails": []string{
						"additional1@example.com",
						"additional2@example.com",
					},
				},
			},
			expected: []CreateFeedback{
				{Email: "additional1@example.com"},
				{Email: "additional2@example.com"},
			},
			hasError: false,
			empty:    false,
		},
		{
			name: "User not manager of the team",
			identity: identity.Identity{
				UID:            uuid.New(),
				OrganizationID: organizationID,
			},
			variables: map[string]interface{}{
				"input": map[string]interface{}{
					"employeeIDs":     []string{emps[0].ID.String(), emps[1].ID.String()},
					"includeTeam":     true,
					"collectionEndAt": time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339),
					"additionalEmails": []string{
						"additional1@example.com",
						"additional2@example.com",
					},
				},
			},
			empty:    true,
			hasError: true,
			errorMsg: "you are not a manager of any team",
		},
		{
			name: "Manager not in the same organization",
			identity: identity.Identity{
				UID:            managerUserID,
				OrganizationID: uuid.New(), // Different organization ID
			},
			variables: map[string]interface{}{
				"input": map[string]interface{}{
					"employeeIDs":     []string{emps[0].ID.String(), emps[1].ID.String()},
					"includeTeam":     true,
					"collectionEndAt": time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339),
					"additionalEmails": []string{
						"additional1@example.com",
						"additional2@example.com",
					},
				},
			},
			empty:    true,
			hasError: true,
			errorMsg: "you are not a manager of any team",
		},
		{
			name: "No employees found",
			identity: identity.Identity{
				UID:            managerUserID,
				OrganizationID: organizationID,
			},
			variables: map[string]interface{}{
				"input": map[string]interface{}{
					"employeeIDs":     []string{uuid.New().String(), uuid.New().String()}, // IDs not matching any employees
					"includeTeam":     true,
					"collectionEndAt": time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339),
					"additionalEmails": []string{
						"additional1@example.com",
						"additional2@example.com",
					},
				},
			},
			empty:    true,
			hasError: true,
			errorMsg: "you do not have any employees matching the criteria",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set the identity in the context
			gctx = identity.ToContext(gctx, tc.identity)

			// Define the test mutation with variables
			mutation := `
				mutation createFeedbacks($input: CreateFeedbacksInput!) {
					createFeedbacks(input: $input) {
						id
						email
					}
				}
			`

			r, err := gql.ExecuteGraphQLMutation(gctx, mutations, mutation, tc.variables)
			assert.NotNil(t, r)

			if tc.hasError {
				assert.NotNil(t, r.Errors)
				if tc.errorMsg != "" {
					assert.Contains(t, r.Errors[0].Message, tc.errorMsg)
				}
				return
			}

			assert.NoError(t, err)
			assert.Nil(t, r.Errors)

			b, _ := json.Marshal(r.Data.(map[string]interface{})["createFeedbacks"])

			var results []CreateFeedback

			err = json.Unmarshal(b, &results)
			assert.NoError(t, err)

			if tc.empty {
				assert.Len(t, results, 0)
				return
			}

			for _, expected := range tc.expected {
				result := golly.Find(results, func(result CreateFeedback) bool {
					return strings.EqualFold(expected.Email, result.Email)
				})

				assert.NotNil(t, result)
				assert.Equal(t, result.Email, expected.Email)
			}
		})
	}
}
