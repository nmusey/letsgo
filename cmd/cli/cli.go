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

	if err := NewApp(name, repo, pwd); err != nil {
		log.Fatal("Unable to create new app: " + err.Error())
	}

	fmt.Printf("Created app %s in directory %s\n", name, pwd)
	return nil
}
