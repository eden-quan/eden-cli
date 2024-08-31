package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/eden/eden-cli/project"
)

func CreateModuleCmd(root *cobra.Command) {
	serviceCmd := &cobra.Command{
		Short: "manage service module",
		Use:   "module",
		Run: func(cmd *cobra.Command, _ []string) {
			_ = cmd.Help()
		},
	}

	serviceNewCmd := &cobra.Command{
		Short: "create new module",
		Use:   "new",
		Run: func(cmd *cobra.Command, args []string) {

			proj, err := CreateModule()
			if err != nil {
				cmd.PrintErr(err)
				return
			}

			if !proj.Exists() {
				cmd.PrintErrf("Project %s doesn't exists, please create project first\n", proj.ProjectName)
				return
			}

			// 创建示例服务
			proj.OverwriteAll = 1
			proj.New("template/module")
			InitialService(proj)
			proj.Done()
		},
	}

	ls := &cobra.Command{
		Short: "list all module on service",
		Use:   "ls",
		Run: func(cmd *cobra.Command, args []string) {
			proj, err := NewProject()
			if err != nil {
				cmd.PrintErrf("init project info with error %s\n", err)
				return
			}

			svc, err := chooseService(args, "please choose the service need to show")
			if err != nil {
				cmd.PrintErrf("choose service failed with error %s\n", err)
				return
			}

			modules, err := proj.ScanModules(svc)
			if err != nil {
				cmd.PrintErrf("scan modules failed with error %s\n", err)
				return
			}

			cmd.Print("modules:\n")
			for i, mod := range modules {
				cmd.Printf("  %d. %s\n", i+1, mod)
			}
		},
	}

	serviceCmd.AddCommand(ls)
	serviceCmd.AddCommand(serviceNewCmd)
	root.AddCommand(serviceCmd)
}

func CreateModule() (project.Project, error) {
	proj, err := NewProject()
	if err != nil {
		return proj, err
	}

	p := ""
	if len(proj.ServiceList) > 1 {
		sel := promptui.Select{
			Label: "service",
			Items: proj.ServiceList,
			Templates: &promptui.SelectTemplates{
				Active:   "[{{.Name}}]",
				Inactive: "{{.Name}}",
				Selected: "Selected: {{.Name}}",
			},
		}
		i := 0
		i, _, err = sel.Run()
		proj.ServiceName = proj.ServiceList[i].Name
	} else if len(proj.ServiceList) == 1 {
		proj.ServiceName = proj.ServiceList[0].Name
		fmt.Printf("Seleted: %s (only one service exists)\n", proj.ServiceName)
	} else {
		return proj, errors.New("so service exists, please create service first")
	}

	// version
	prompt := promptui.Prompt{
		Label:   "project version",
		Default: "v1",
	}
	p, err = prompt.Run()
	proj.Version = p

	if err == nil {
		defaultModule := strings.ReplaceAll(proj.ServiceName, "-service", "")
		// check if exists
		moduleList, _ := proj.ScanModules(proj.ServiceName)

		// module
		prompt := promptui.Prompt{
			Label:    "module",
			Default:  defaultModule,
			Validate: Validator(moduleList, "module"),
		}
		p, err = prompt.Run()
		proj.ModuleName = p
	}

	proj.Prepare()

	return proj, err
}
