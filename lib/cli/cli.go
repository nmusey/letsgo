package cli

import (
	"context"

	"github.com/nmusey/letsgo/lib/templates"
	"github.com/urfave/cli/v3"
)

func MakeCommands() *cli.Command {
	return &cli.Command{
		Commands: []*cli.Command{
			{
				Name:    "make",
				Aliases: []string{"m"},
				Commands: []*cli.Command{
					{
						Name: "project",
						Usage: "Scaffold a go project quickly",
						Action: func(context.Context, *cli.Command) error {
							template := templates.NewTestTemplate("substitution")
							return template.ParseFiles()
						},
					},
				},
			},
		},
	}
}
