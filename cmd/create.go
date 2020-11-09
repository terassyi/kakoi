package cmd

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"github.com/terassyi/kakoi/infra"
)

type Create struct {
	Path string
}

func (*Create) Name() string {
	return "create"
}

func (*Create) Usage() string {
	return "kakoi create -p [file path]"
}

func (*Create) Synopsis() string {
	return "create infrastructure resources"
}

func (c *Create) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.Path, "p", "", "config file path")
}

func (c *Create) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	i, err := infra.New(c.Path)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}
	if err := i.Create(); err != nil {
		fmt.Println(err)
	}
	return subcommands.ExitSuccess
}
