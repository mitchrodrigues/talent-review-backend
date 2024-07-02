package cycle

import (
	"testing"
	"time"

	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAggregateApply_CycleCreated(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		event    eventsource.Event
		expected Aggregate
	}{
		{
			name: "CycleCreated event",
			event: eventsource.Event{
				Data: CycleCreated{
					ID:             uuid.New(),
					OrganizationID: uuid.New(),
					OwnerID:        uuid.New(),
					StartAt:        time.Now(),
					EndAt:          time.Now().Add(24 * time.Hour),
				},
				CreatedAt: time.Now(),
			},
			expected: Aggregate{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize aggregate
			aggregate := &Aggregate{}

			// Call the Apply method
			aggregate.Apply(golly.Context{}, tt.event)

			// Validate results
			assert.Equal(t, tt.event.Data.(CycleCreated).ID, aggregate.ID)
			assert.Equal(t, tt.event.Data.(CycleCreated).OrganizationID, aggregate.OrganizationID)
			assert.Equal(t, tt.event.Data.(CycleCreated).OwnerID, aggregate.OwnerID)
			assert.Equal(t, tt.event.Data.(CycleCreated).StartAt, aggregate.StartAt)
			assert.Equal(t, tt.event.Data.(CycleCreated).EndAt, aggregate.EndAt)
			assert.Equal(t, tt.event.CreatedAt, aggregate.CreatedAt)
			assert.Equal(t, tt.event.CreatedAt, aggregate.UpdatedAt)
		})
	}
}
