package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/eden-quan/eden-cli/project"
)

func CreateProjCmd(root *cobra.Command) {
	projCmd := &cobra.Command{
		Short: "project manager",
		Use:   "project",
		Run: func(cmd *cobra.Command, _ []string) {
			_ = cmd.Help()
		},
	}

	projNewCmd := &cobra.Command{
		Short: "create new project",
		Use:   "new",
		Run: func(cmd *cobra.Command, args []string) {

			proj, err := CreateProject()
			if err != nil {
				cmd.PrintErr(err)
				return
			}

			proj.OverwriteAll = 1
			// 创建项目
			proj.New("template/project")
			proj.New("template/project-upgrade")

			// 创建示例服务
			proj.New("template/service")

			// 创建示例模块
			proj.New("template/module")

			InitialProject(proj)
			proj.Done()

		},
	}

	projUpgradeCmd := &cobra.Command{
		Short: "upgrade current project",
		Use:   "upgrade",
		Run: func(cmd *cobra.Command, args []string) {

			proj, err := UpgradeProject()
			if err != nil {
				cmd.PrintErr(err)
				return
			}

			// 创建项目
			proj.New("template/project-upgrade")

			InitialProject(proj)
			proj.Done()

		},
	}

	projCmd.AddCommand(projNewCmd)
	projCmd.AddCommand(projUpgradeCmd)
	root.AddCommand(projCmd)
}

func CreateProject() (project.Project, error) {
	// dir
	basePath, _ := filepath.Abs(".")
	proj := project.Project{}
	prompt := promptui.Prompt{
		Label:    "project path",
		Default:  basePath,
		Validate: Validator(nil, "project"),
	}
	p, err := prompt.Run()

	proj.BasePath = p

	if err == nil {
		// project name
		prompt = promptui.Prompt{
			Label:   "project name",
			Default: "some-service",
		}
		p, err = prompt.Run()
		proj.ProjectName = p
		proj.CheckProjectName()
	}

	if err == nil {
		// package name
		prompt = promptui.Prompt{
			Label:   "project package",
			Default: fmt.Sprintf("github.com/eden-quan/%s", proj.ProjectName),
		}
		p, err = prompt.Run()
		proj.PackageName = strings.TrimSpace(p)
	}

	if err == nil {
		// version
		prompt = promptui.Prompt{
			Label:   "project version",
			Default: "v1",
		}
		p, err = prompt.Run()
		proj.Version = p
	}

	if err == nil {
		// service
		prompt = promptui.Prompt{
			Label:   "default service",
			Default: proj.ProjectName,
		}
		p, err = prompt.Run()
		proj.ServiceName = p
	}

	if err == nil {
		// module
		prompt = promptui.Prompt{
			Label:   "default module",
			Default: strings.Replace(proj.ProjectName, "-service", "", 1),
		}
		p, err = prompt.Run()
		proj.ModuleName = p
	}

	proj.Prepare()

	return proj, err

}

func UpgradeProject() (project.Project, error) {
	proj, err := NewProject()
	if err != nil {
		return proj, err
	}

	// version
	prompt := promptui.Prompt{
		Label:   "project version",
		Default: "v1",
	}
	p, err := prompt.Run()
	proj.Version = p
	proj.Prepare()

	selectPrompt := promptui.Select{
		Label: "are you sure for upgrade? it will overwrite all the basic file (excluded services)",
		Items: []string{"Y", "N"},
	}
	n, _, err := selectPrompt.Run()

	if n == 1 || err != nil {
		return proj, errors.New("cancel")
	}

	return proj, err
}

func InitialProject(project project.Project) {
	cmdList := []*exec.Cmd{
		exec.Command("mv", "gitignore", ".gitignore"),
		exec.Command("make", "init"),
		exec.Command("make", "proto-gen", fmt.Sprintf("service=%s", project.ServiceShortName)),
		exec.Command("goimports", "-w", "api"),
		exec.Command("go", "mod", "tidy"),
	}

	for _, cmd := range cmdList {
		cmd.Dir = project.BasePath
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("calling cmd %s with error %s", cmd.String(), err)
			os.Exit(0)
		}
	}
}

func UpgradeCli() {
	cmdList := []*exec.Cmd{
		exec.Command("go", "install", "github.com/eden-quan/eden-cli@latest"),
	}

	for _, cmd := range cmdList {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("upgrade eden-cli with cmd %s occur error %s", cmd.String(), err)
			os.Exit(0)
		}
	}
}
