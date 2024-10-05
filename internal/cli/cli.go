package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)


func RunCli() {
	app := cli.App {
		Name:    "letsgo",
		Usage:   "letsgo --help",
		Suggest: true,
		Commands: []*cli.Command{
			{
				Name:  "make",
				Usage: "letsgo make name repo",
                Action: makeApp,
			},
            {
                Name: "pkg",
                Usage: "letsgo pkg name",
                Action: makePackage,
            },
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func makeApp(ctx *cli.Context) error {
	name := ctx.Args().Get(0)
    repo := ctx.Args().Get(1)
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := NewAppTemplate(name, repo, pwd); err != nil {
		log.Fatal("Unable to create new app: " + err.Error())
	}

	fmt.Printf("Created app %s in directory %s\n", name, pwd)
	return nil
}

func makePackage(ctx *cli.Context) error {
    name := ctx.Args().Get(0)
    pwd, err := os.Getwd()
    if err != nil {
        return err
    }

    if err := NewPackageTemplate(name, pwd); err != nil {
        log.Fatal("Unable to create package: " + err.Error())
    }

    return nil
}
