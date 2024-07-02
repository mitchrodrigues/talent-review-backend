package cycle

import (
	"context"
	"testing"
	"time"

	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/orm"
	"github.com/google/uuid"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/identity"
	"github.com/stretchr/testify/assert"
)

func TestFindOrCreateCycle_Perform(t *testing.T) {
	// Initialize in-memory database
	gctx := orm.CreateTestContext(golly.NewContext(context.TODO()), Aggregate{})
	ident, gctx := identity.NewTestIdentity(gctx)

	existingID := uuid.New()

	// Define test cases
	tests := []struct {
		name       string
		cmd        FindOrCreateCycle
		setup      func()
		expectErr  bool
		shouldFind bool
	}{
		{
			name: "Create new cycle",
			cmd: FindOrCreateCycle{
				StartAt: time.Now(),
				EndAt:   time.Now().Add(24 * time.Hour),
				Type:    "type1",
			},
			setup:      func() {},
			expectErr:  false,
			shouldFind: false,
		},
		{
			name: "Find existing cycle",
			cmd: FindOrCreateCycle{
				StartAt: time.Now().Add(-12 * time.Hour),
				EndAt:   time.Now().Add(12 * time.Hour),
				Type:    "type2",
			},
			setup: func() {
				orm.DB(gctx).Create(&Aggregate{
					ModelUUID: orm.ModelUUID{
						ID: existingID,
					},
					OrganizationID: ident.OrganizationID,
					OwnerID:        ident.UID,
					StartAt:        time.Now().Add(-24 * time.Hour),
					EndAt:          time.Now().Add(24 * time.Hour),
				})
			},
			expectErr:  false,
			shouldFind: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test case
			tt.setup()

			// Initialize aggregate
			aggregate := &Aggregate{}

			// Call the Perform method
			err := tt.cmd.Perform(gctx, aggregate)

			// Validate results
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.shouldFind {
				assert.Equal(t, existingID, aggregate.ID)
			} else {
				assert.NotEqual(t, existingID, aggregate.ID)
			}
		})
	}
}
