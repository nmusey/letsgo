package cli

import (
	"context"

	"github.com/nmusey/letsgo/lib/templates"
	"github.com/urfave/cli/v3"
)

func MakeCommands() *cli.Command {
	return &cli.Command{
		Name:  "letsgo",
		Usage: "Scaffold a go project quickly",
		Action: func(context.Context, *cli.Command) error {
			template := templates.NewTestTemplate("substitution")
			return template.ParseFiles()
		},
	}
}
