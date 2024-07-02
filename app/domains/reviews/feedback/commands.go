package feedback

import (
	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
	"github.com/google/uuid"
)

type CreateFeedback struct {
	TeamID uuid.UUID
}

func (cmd CreateFeedback) Perform(gctx golly.Context, aggregate eventsource.Aggregate) error {
	return nil
}
