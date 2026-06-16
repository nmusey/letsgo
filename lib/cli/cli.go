package cli

import (
	"context"

	"github.com/nmusey/letsgo/lib/files"
	"github.com/nmusey/letsgo/lib/templates"
	"github.com/urfave/cli/v3"
)

func MakeCommands() *cli.Command {
	return &cli.Command{
		Name:  "letsgo",
		Usage: "Scaffold a go project quickly",
		Action: func(context.Context, *cli.Command) error {
			substitutions := files.Substitutions{
				AppName: "temp",
				DbPort: 5432,
			}

			return files.ParseFiles(templates.ProjectTemplate, substitutions)
		},
	}
}
