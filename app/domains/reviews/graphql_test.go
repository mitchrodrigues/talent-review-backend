package reviews

import (
	"context"
	"encoding/json"
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
	}

	// Set up the test context and seed the database
	gctx := createTestContext()

	// User ID for the manager
	managerUserID := uuid.New()
	organizationID := uuid.New()

	manager := employees.NewTestEmployee(uuid.New(), organizationID, &managerUserID)
	orm.DB(gctx).Create(&manager)

	// Team with the manager
	team := employees.NewTestTeam(uuid.New(), organizationID, manager.ID)
	orm.DB(gctx).Create(&team)

	// Employees in the team
	emps := make([]employees.Employee, 2)
	for pos := range emps {
		emps[pos] = employees.NewTestEmployee(uuid.New(), organizationID, nil)
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
			hasError: true,
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
			hasError: true,
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
				return
			}

			assert.NoError(t, err)
			assert.Nil(t, r.Errors)

			b, _ := json.Marshal(r.Data.(map[string]interface{})["createFeedbacks"])

			var results []CreateFeedback

			err = json.Unmarshal(b, &results)
			assert.NoError(t, err)

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
