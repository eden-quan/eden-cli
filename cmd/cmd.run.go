package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func CreateRunAndGen(root *cobra.Command) {
	run := &cobra.Command{
		Short: "run the target service",
		Use:   "run",
		Run: func(cmd *cobra.Command, args []string) {
			service, err := chooseService(args, "please choose the service need to run")
			if err != nil {
				cmd.PrintErr(err)
			}

			c := exec.Command("make", "run", fmt.Sprintf("service=%s", service))
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr

			err = c.Run()
			if err != nil {
				cmd.PrintErrf("run service %s with error %s\n", service, err)
			}
		},
	}

	gen := &cobra.Command{
		Short: "generate codes by eden plugins",
		Use:   "gen",
		Run: func(cmd *cobra.Command, args []string) {
			service, err := chooseService(args, "please choose the service need to gen")
			if err != nil {
				cmd.PrintErr(err)
			}

			c := exec.Command("make", "proto-gen", fmt.Sprintf("service=%s", service))
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			err = c.Run()

			if err == nil {
				c = exec.Command("goimports", "-w", fmt.Sprintf("api/%s", service))
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				err = c.Run()
			}

			if err != nil {
				cmd.PrintErrf("generate service %s with error %s\n", service, err)
			}
		},
	}

	root.AddCommand(gen)
	root.AddCommand(run)

}
