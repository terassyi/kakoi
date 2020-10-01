package cmd

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"github.com/terassyi/kakoi/infra"
	"github.com/terassyi/kakoi/infra/terraform"
)

type Destroy struct {
	Path string
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
}

func (d *Destroy) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	workDir, err := infra.IsExistWorkDir(d.Path)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		return subcommands.ExitFailure
	}
	fmt.Println("[INFO] destroy")
	tf, err := terraform.Prepare(workDir)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		return subcommands.ExitFailure
	}
	tfCtx := context.Background()
	if err := tf.Destroy(tfCtx); err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
