package main

import (
	"context"
	"flag"
	"github.com/google/subcommands"
	"github.com/terassyi/kakoi/cmd"
	"os"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")

	subcommands.Register(&cmd.Init{}, "internal")
	subcommands.Register(&cmd.Create{}, "infra")
	subcommands.Register(&cmd.Destroy{}, "infra")


	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
