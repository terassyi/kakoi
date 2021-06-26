package cmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
	"github.com/terassyi/kakoi/infra"
)

type Destroy struct {
	Path  string
	Force bool
}

func (*Destroy) Name() string {
	return "destroy"
}

func (*Destroy) Synopsis() string {
	return "kakoi destroy -p [path/to/config] destroy infra resource"
}

func (*Destroy) Usage() string {
	return "kakoi destroy -p [path/to/config] destroy infra resource"
}

func (d *Destroy) SetFlags(f *flag.FlagSet) {
	f.StringVar(&d.Path, "p", "", "destroy resource")
	f.BoolVar(&d.Force, "f", false, "force destroy")
}

func (d *Destroy) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	destroyer, err := infra.NewDestroyer(d.Path, d.Force)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		return subcommands.ExitFailure
	}
	if err := destroyer.Destroy(); err != nil {
		fmt.Println("[ERROR] ", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
