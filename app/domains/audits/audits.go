package audits

import "github.com/golly-go/golly"

func Initialize(app golly.Application) error {
	InitGraphQL()

	return nil
}
