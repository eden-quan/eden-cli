package main

import "github.com/spf13/cobra"

const Version = "0.0.1"

func VersionCmd() *cobra.Command {

	versionCmd := &cobra.Command{
		Short: "version of this application",
		Use:   "version",
		Run: func(cmd *cobra.Command, _ []string) {
			_ = cmd.Help()
			cmd.Printf("eden-cli %s\n", Version)
		},
	}

	return versionCmd
}
