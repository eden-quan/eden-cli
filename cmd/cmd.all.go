package cmd

import (
	"slices"

	"github.com/spf13/cobra"

	"github.com/eden-quan/eden-cli/project"
)

func CreateAllInOneCmd(root *cobra.Command) {
	serviceCmd := &cobra.Command{
		Short: "merge all services in this project to an single application",
		Long:  "merge all services in this project to an service named: all-in-one-service, every service will use an independence environment",
		Use:   "all-in-one",
		Run: func(cmd *cobra.Command, args []string) {
			proj, err := NewProject()
			if err != nil {
				cmd.PrintErrf("creating all in one mode failed with error %s\n", err)
				return
			}

			proj.Prepare()
			// remove all-in-one
			proj.ServiceList = slices.DeleteFunc(proj.ServiceList, func(name project.ServiceName) bool {
				return name.Name == "all-in-one-service"
			})
			proj.New("template/all-in-one")
			proj.OverwriteAll = 1
			InitialService(proj)
			proj.Done()
		},
	}

	root.AddCommand(serviceCmd)
}
