package cmd

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"github.com/terassyi/kakoi/infra"
)

type Init struct {
	Path string
}

func (*Init) Name() string {
	return "init"
}

func (*Init) Usage() string {
	return "kakoi init -p [path/to/config/file]"
}

func (*Init) Synopsis() string {
	return "initialize kakoi state and resources"
}

func (i *Init) SetFlags(f *flag.FlagSet) {
	f.StringVar(&i.Path, "p", "", "set path to config file.")
}

func (i *Init) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	initer, err := infra.NewInitializer(i.Path)
	if err != nil {
		fmt.Println("[ERROR] ", err)
		return subcommands.ExitFailure
	}
	if err := initer.Init(); err != nil {
		fmt.Println("[ERROR] ", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
