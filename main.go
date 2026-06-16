package main

import (
	"context"
	"log"
	"os"

	"github.com/nmusey/letsgo/lib/cli"
)

func main() {
	cmd := cli.MakeCommands()
	if err := cmd.Run(context.Background(), os.Args); err != nil {
        log.Fatal(err)
    }
}
