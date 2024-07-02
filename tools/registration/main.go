package main

import (
	"time"

	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/eventsource"
	"github.com/golly-go/plugins/orm"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/accounts"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/accounts/organizations"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/accounts/users"
	"github.com/mitchrodrigues/talent-review-backend/app/initializers"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/workos"
	"github.com/spf13/cobra"
)

var commands = []*cobra.Command{
	{
		Use:  "register",
		Long: "create an organization and user",
		Args: cobra.MinimumNArgs(4),
		Run:  golly.Command(register),
	},
	{
		Use:  "invite [orgID] [email]",
		Long: "create an organization and user",
		Args: cobra.MinimumNArgs(2),
		Run:  golly.Command(invite),
	},
}

func main() {
	golly.Start(golly.GollyStartOptions{
		Preboots:     initializers.Preboots,
		Initializers: initializers.Initializers,
		CLICommands:  commands,
	})
}

func invite(gctx golly.Context, cmd *cobra.Command, args []string) error {
	var org organizations.Aggregate

	if err := orm.DB(gctx).Model(org).Find(&org, "id = ?", args[0]).Error; err != nil {
		return err
	}

	user := accounts.User{}

	err := eventsource.Call(gctx, &user.Aggregate, users.InviteUser{
		WorkosClient: workos.DefaultClient{},
		Email:        args[1],
		Organization: &org,
	}, eventsource.Metadata{})

	if err != nil {
		return err
	}

	time.Sleep(1 * time.Second)
	return nil
}

func register(gctx golly.Context, cmd *cobra.Command, args []string) error {
	organization := organizations.Aggregate{}

	{
		err := eventsource.Call(gctx, &organization, organizations.CreateOrganization{
			WorkosClient: workos.DefaultClient{},
			Name:         args[0],
		}, eventsource.Metadata{})

		if err != nil {
			return err
		}
	}

	user := accounts.User{}

	return eventsource.Call(gctx, &user.Aggregate, users.CreateUser{
		WorkosClient: workos.DefaultClient{},
		FirstName:    args[1],
		LastName:     args[2],
		Email:        args[3],
		Organization: &organization,
	}, eventsource.Metadata{})
}
