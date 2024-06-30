package esbackend

import (
	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
)

type Backend struct {
	PostgresRepository
}

func (b Backend) PublishEvent(ctx golly.Context, aggregate eventsource.Aggregate, data ...eventsource.Event) {
}

func Initializer(app golly.Application) error {
	eventsource.SetEventRepository(Backend{})
	return nil
}
