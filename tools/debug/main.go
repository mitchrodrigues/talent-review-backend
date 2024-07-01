package main

import (
	"fmt"

	"github.com/golly-go/golly"
	"github.com/mitchrodrigues/talent-review-backend/app/initializers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var commands = []*cobra.Command{
	{
		Use:  "env",
		Long: "print-env",
		Run:  golly.Command(printEnv),
	},
}

func main() {
	golly.Start(golly.GollyStartOptions{
		Preboots:     initializers.Preboots,
		Initializers: initializers.Initializers,
		CLICommands:  commands,
	})
}

func printEnv(gctx golly.Context, cmd *cobra.Command, args []string) error {

	keys := gctx.Config().AllKeys()
	for _, key := range keys {
		fmt.Printf("%s: %v", key, viper.Get(key))
	}

	return nil
}
