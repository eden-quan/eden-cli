package main

import (
	"github.com/spf13/cobra"

	"gitlab.lainuoniao.cn/rhinobird/backend/eden-cli.git/cmd"
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
