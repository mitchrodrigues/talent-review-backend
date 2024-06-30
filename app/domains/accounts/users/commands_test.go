package users

import (
	"errors"
	"testing"

	"github.com/golly-go/golly"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/helpers"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/workos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockIDPObject is a mock implementation of the IDPObject interface
type MockIDPObject struct {
	idpID    string
	recordID uuid.UUID
}

func (m *MockIDPObject) RecordIdpID() string {
	return m.idpID
}

func (m *MockIDPObject) RecordID() uuid.UUID {
	return m.recordID
}
func TestCreateUserPerform(t *testing.T) {
	tests := []struct {
		name      string
		cmd       CreateUser
		expectErr bool
		mockSetup func() *workos.MockClient
	}{
		{
			name: "success",
			cmd: CreateUser{
				Organization: &MockIDPObject{
					idpID:    "org-id",
					recordID: uuid.New(),
				},
				Email:     "test@example.com",
				FirstName: "Test",
				LastName:  "User",
				Password:  "password",
			},
			expectErr: false,
			mockSetup: func() *workos.MockClient {
				mockClient := &workos.MockClient{}

				mockClient.On("CreateUser", mock.Anything, workos.CreateUserInput{
					FirstName:     "Test",
					LastName:      "User",
					Email:         "test@example.com",
					Password:      "password",
					OrganzationID: "org-id",
				}).Return("idp-id", nil)
				return mockClient
			},
		},
		{
			name: "failure in CreateUser",
			cmd: CreateUser{
				Organization: &MockIDPObject{
					idpID:    "org-id",
					recordID: uuid.New(),
				},
				Email:     "test@example.com",
				FirstName: "Test",
				LastName:  "User",
				Password:  "password",
			},
			expectErr: true,
			mockSetup: func() *workos.MockClient {
				mockClient := &workos.MockClient{}
				mockClient.On("CreateUser", mock.Anything, workos.CreateUserInput{
					FirstName:     "Test",
					LastName:      "User",
					Email:         "test@example.com",
					Password:      "password",
					OrganzationID: "org-id",
				}).Return("", errors.New("failed to create user"))
				return mockClient
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := golly.Context{}
			user := Aggregate{}

			mock := tt.mockSetup()
			tt.cmd.WorkosClient = mock

			err := tt.cmd.Perform(ctx, &user)
			if tt.expectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			changes := user.Changes()
			assert.Len(t, changes, 1)

			event, ok := changes[0].Data.(UserCreated)
			assert.True(t, ok)

			assert.Equal(t, tt.cmd.FirstName, event.FirstName)
			assert.Equal(t, tt.cmd.LastName, event.LastName)
			assert.Equal(t, tt.cmd.Email, event.Email)
			assert.Equal(t, tt.cmd.Organization.RecordID(), event.OrganizationID)

			mock.AssertExpectations(t)
		})
	}
}

func TestInviteUserValidate(t *testing.T) {
	tests := []struct {
		name      string
		cmd       InviteUser
		expectErr bool
	}{
		{
			name: "missing organization",
			cmd: InviteUser{
				WorkosClient: &workos.MockClient{},
				Email:        "test@example.com",
			},
			expectErr: true,
		},
		{
			name: "missing workos client",
			cmd: InviteUser{
				Organization: &MockIDPObject{
					idpID:    "org-id",
					recordID: uuid.New(),
				},
				Email: "test@example.com",
			},
			expectErr: true,
		},
		{
			name: "valid input",
			cmd: InviteUser{
				WorkosClient: &workos.MockClient{},
				Organization: &MockIDPObject{
					idpID:    "org-id",
					recordID: uuid.New(),
				},
				Email: "test@example.com",
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := golly.Context{}
			user := Aggregate{}

			err := tt.cmd.Validate(ctx, &user)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestInviteUserPerform(t *testing.T) {
	tests := []struct {
		name      string
		cmd       InviteUser
		expectErr bool
		mockSetup func() *workos.MockClient
	}{
		{
			name: "success",
			cmd: InviteUser{
				Organization: &MockIDPObject{
					idpID:    "org-id",
					recordID: uuid.New(),
				},
				Inviter: &MockIDPObject{
					idpID:    "inviter-id",
					recordID: uuid.New(),
				},
				Email: "test@example.com",
			},
			expectErr: false,
			mockSetup: func() *workos.MockClient {
				mockClient := &workos.MockClient{}
				mockClient.On("InviteUser", mock.Anything, "org-id", "test@example.com", "inviter-id").Return("idp-id", nil)
				return mockClient
			},
		},
		{
			name: "failure in InviteUser",
			cmd: InviteUser{
				Organization: &MockIDPObject{
					idpID:    "org-id",
					recordID: uuid.New(),
				},
				Inviter: &MockIDPObject{
					idpID:    "inviter-id",
					recordID: uuid.New(),
				},
				Email: "test@example.com",
			},
			expectErr: true,
			mockSetup: func() *workos.MockClient {
				mockClient := &workos.MockClient{}
				mockClient.On("InviteUser", mock.Anything, "org-id", "test@example.com", "inviter-id").Return("", errors.New("failed to invite user"))
				return mockClient
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := golly.Context{}
			user := Aggregate{}

			mock := tt.mockSetup()
			tt.cmd.WorkosClient = mock

			err := tt.cmd.Perform(ctx, &user)
			if tt.expectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			mock.AssertExpectations(t)

			changes := user.Changes()
			assert.Len(t, changes, 1)

			event, ok := changes[0].Data.(UserCreated)
			assert.True(t, ok)

			assert.Equal(t, tt.cmd.Email, event.Email)
			// assert.Equal(t, tt.expectedEvent.IdpInviteID, event.IdpInviteID)
			assert.Equal(t, tt.cmd.Organization.RecordID(), event.OrganizationID)
			if tt.cmd.Inviter != nil {
				assert.Equal(t, tt.cmd.Inviter.RecordID(), event.InviterID)
			}
		})
	}
}

func TestEditUserPerform(t *testing.T) {
	tests := []struct {
		name      string
		cmd       EditUser
		existing  *Aggregate
		expectErr bool
	}{
		{
			name: "success",
			cmd: EditUser{
				IdpID:     "new-idp-id",
				FirstName: "NewFirstName",
				LastName:  "NewLastName",
				Email:     "new@example.com",
			},
			existing: &Aggregate{
				IdpID:     "old-idp-id",
				FirstName: "OldFirstName",
				LastName:  "OldLastName",
				Email:     "old@example.com",
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := golly.Context{}

			aggregate := tt.existing

			err := tt.cmd.Perform(ctx, aggregate)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				changes := aggregate.Changes()
				assert.Len(t, changes, 1)
				event, ok := changes[0].Data.(UserUpdated)
				assert.True(t, ok)

				assert.Equal(t, helpers.Coalesce(tt.cmd.IdpID, tt.existing.IdpID), event.IdpID)
				assert.Equal(t, helpers.Coalesce(tt.cmd.FirstName, tt.existing.FirstName), event.FirstName)
				assert.Equal(t, helpers.Coalesce(tt.cmd.LastName, tt.existing.LastName), event.LastName)
				assert.Equal(t, helpers.Coalesce(tt.cmd.Email, tt.existing.Email), event.Email)
			}
		})
	}
}
