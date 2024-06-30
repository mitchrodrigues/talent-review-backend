package main

import (
	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/orm/migrate"
	"github.com/mitchrodrigues/talent-review-backend/app/initializers"
)

var commands = append(golly.AppCommands, migrate.Command())

func main() {
	golly.Start(golly.GollyStartOptions{
		Preboots:     initializers.Preboots,
		Initializers: initializers.Initializers,
		CLICommands:  commands,
	})
}
