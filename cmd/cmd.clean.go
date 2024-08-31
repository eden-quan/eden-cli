package cmd

import (
	"github.com/spf13/cobra"
)

func CreateCleanCmd(root *cobra.Command) {
	serviceCmd := &cobra.Command{
		Short: "clean all the generated file",
		Long:  "clean all the generated file from protocol buffer definition",
		Use:   "clean",
		Run: func(cmd *cobra.Command, args []string) {
			proj, err := NewProject()
			if err != nil {
				cmd.PrintErrf("creating clean command failed with error %s\n", err)
				return
			}

			proj.Prepare()
			CleanService(proj)
			proj.Done()
		},
	}

	root.AddCommand(serviceCmd)
}

func CreateUpgrade(root *cobra.Command) {
	serviceCmd := &cobra.Command{
		Short: "upgrade eden-cli",
		Long:  "upgrade eden-cli to latest version",
		Use:   "upgrade",
		Run: func(cmd *cobra.Command, args []string) {
			UpgradeCli()
		},
	}

	root.AddCommand(serviceCmd)
}
