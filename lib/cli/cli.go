package cli

import (
	"context"
	"path"

	"github.com/nmusey/letsgo/lib/names"
	"github.com/nmusey/letsgo/lib/templates"
	"github.com/urfave/cli/v3"
)

func MakeCommands() *cli.Command {
	makeCommands := []*cli.Command{
		{
			Name:  "project",
			Usage: "Scaffold a go project quickly",
			Arguments: []cli.Argument{
				&cli.StringArg{
					Name:     "name",
					Required: true,
				},
				&cli.StringArg{
					Name:     "repository",
					Required: true,
				},
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				name := cmd.StringArg("name")
				repository := cmd.StringArg("repository")
				template := templates.NewProjectTemplate(name, repository)
				return template.ParseFiles(name)
			},
		},
		{
			Name:  "model",
			Usage: "Add a domain model to the project. Includes a repository and service alongside the model struct itself.",
			Arguments: []cli.Argument{
				&cli.StringArg{
					Name:     "name",
					Required: true,
				},
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				name := cmd.StringArg("name")
				template, err := templates.NewModelTemplate(name)
				if err != nil {
					return err
				}

				return template.ParseFiles(path.Join("lib", "domain", names.PackageName(name)))
			},
		},
		{
			Name:  "app-service",
			Usage: "Add an app service to the project.",
			Arguments: []cli.Argument{
				&cli.StringArg{
					Name:     "name",
					Required: true,
				},
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				name := cmd.StringArg("name")
				template, err := templates.NewAppServiceTemplate(name)
				if err != nil {
					return err
				}

				return template.ParseFiles(path.Join("lib", "domain", names.PackageName(name)))
			},
		},
	}

	return &cli.Command{
		Commands: []*cli.Command{
			{
				Name:     "make",
				Aliases:  []string{"m"},
				Commands: makeCommands,
			},
		},
	}
}
