package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/nmusey/letsgo/lib/generate"
	"github.com/nmusey/letsgo/lib/templates"
	"github.com/urfave/cli/v3"
)

func MakeCommands() *cli.Command {
	makeCommands := []*cli.Command{
		{
			Name:  "project",
			Usage: "Scaffold a go project quickly",
			Action: func(context.Context, *cli.Command) error {
				template := templates.NewTestTemplate("substitution")
				return template.ParseFiles()
			},
		},
	}
	makeCommands = append(makeCommands, generatorCommands()...)

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

func generatorCommands() []*cli.Command {
	commands := make([]*cli.Command, 0, len(generate.Kinds))

	for _, kind := range generate.Kinds {
		commands = append(commands, generatorCommand(kind))
	}

	return commands
}

func generatorCommand(kind generate.Kind) *cli.Command {
	var companions []generate.Kind
	for _, other := range generate.Kinds {
		if other.Name != kind.Name {
			companions = append(companions, other)
		}
	}

	flags := make([]cli.Flag, 0, len(companions))
	for _, companion := range companions {
		flags = append(flags, &cli.BoolFlag{
			Name:    companion.FlagLong,
			Aliases: []string{companion.FlagShort},
			Usage:   fmt.Sprintf("Also generate a %s", companion.Name),
		})
	}

	return &cli.Command{
		Name:  kind.Name,
		Usage: fmt.Sprintf("Generate a %s for a domain", kind.Name),
		Flags: flags,
		Action: func(_ context.Context, cmd *cli.Command) error {
			rawName := cmd.Args().First()
			if rawName == "" {
				return fmt.Errorf("make %s: a domain name is required", kind.Name)
			}

			kindNames := []string{kind.Name}
			for _, companion := range companions {
				if cmd.Bool(companion.FlagLong) {
					kindNames = append(kindNames, companion.Name)
				}
			}

			dir, err := os.Getwd()
			if err != nil {
				return err
			}

			return generate.Run(dir, kindNames, rawName)
		},
	}
}
