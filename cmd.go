package main

import (
	"github.com/spf13/cobra"

	"github.com/eden-quan/eden-cli/cmd"
)

func createCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "eden-cli",
		Annotations: map[string]string{
			cobra.CommandDisplayNameAnnotation: "eden-cli",
		},
	}

	cmd.CreateProjCmd(rootCmd)
	cmd.CreateServiceCmd(rootCmd)
	cmd.CreateModuleCmd(rootCmd)
	cmd.CreateRunAndGen(rootCmd)
	cmd.CreateAllInOneCmd(rootCmd)
	cmd.CreateCleanCmd(rootCmd)
	cmd.CreateUpgrade(rootCmd)

	rootCmd.AddCommand(VersionCmd())

	return rootCmd
}
