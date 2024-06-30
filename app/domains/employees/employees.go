package employees

import (
	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
)

func Initalizer(app golly.Application) error {
	InitGraphQL()

	eventsource.Subscribe("users.Aggregate", "users.UserCreated", UpdateEmployeeUser)

	return nil
}
