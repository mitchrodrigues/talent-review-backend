package employee

import (
	"context"
	"testing"

	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/orm"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/accounts"
	"github.com/stretchr/testify/assert"
)

// Test for Perform method
func TestCreateEmployeePerform(t *testing.T) {
	ctx := orm.CreateTestContext(golly.NewContext(context.TODO()), accounts.User{})

	tests := []struct {
		name      string
		cmd       Create
		expectErr bool
		expected  Created
	}{
		{
			name: "create IC employee",
			cmd: Create{
				Name:           "John Doe",
				Email:          "john.doe@example.com",
				OrganizationID: uuid.New(),
				WorkerType:     "full-time",
			},
			expectErr: false,
			expected: Created{
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
		},
		{
			name: "create Manager employee",
			cmd: Create{
				Name:           "Jane Smith",
				Email:          "jane.smith@example.com",
				OrganizationID: uuid.New(),
				WorkerType:     "part-time",
			},
			expectErr: false,
			expected: Created{
				Name:  "Jane Smith",
				Email: "jane.smith@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			aggregate := &Aggregate{}

			err := tt.cmd.Perform(ctx, aggregate)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				changes := aggregate.Changes()
				assert.GreaterOrEqual(t, len(changes), 1)

				event, ok := changes[0].Data.(Created)

				assert.True(t, ok)
				assert.Equal(t, tt.cmd.Name, event.Name)
				assert.Equal(t, tt.cmd.Email, event.Email)
				assert.Equal(t, tt.cmd.OrganizationID, event.OrganizationID)
			}
		})
	}
}
