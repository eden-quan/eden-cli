package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/eden/eden-cli/project"
)

func CreateServiceCmd(root *cobra.Command) {
	serviceCmd := &cobra.Command{
		Short: "service manager",
		Use:   "service",
		Run: func(cmd *cobra.Command, _ []string) {
			_ = cmd.Help()
		},
	}

	serviceNewCmd := &cobra.Command{
		Short: "create new service",
		Use:   "new",
		Run: func(cmd *cobra.Command, args []string) {

			proj, err := CreateService()
			if err != nil {
				cmd.PrintErr(err)
				return
			}

			if !proj.Exists() {
				cmd.PrintErrf("Project %s doesn't exists, please create project first\n", proj.ProjectName)
				return
			}

			// 创建示例服务
			proj.New("template/service")
			// 创建示例模块
			proj.New("template/module")

			proj.OverwriteAll = 1
			InitialService(proj)
			proj.Done()
		},
	}

	ls := &cobra.Command{
		Short: "list services",
		Use:   "ls",
		Run: func(cmd *cobra.Command, args []string) {
			services, err := os.ReadDir("./api")
			if err != nil {
				cmd.PrintErrf("list service failed with error %s\n", err)
			}

			names := make([]string, 0)
			for _, svc := range services {
				if strings.Index(svc.Name(), "-service") >= 0 {
					names = append(names, svc.Name())
				}
			}

			cmd.Print("services:\n")
			for i, svc := range names {
				cmd.Printf("  %d. %s\n", i+1, svc)
			}
		},
	}

	serviceCmd.AddCommand(ls)
	serviceCmd.AddCommand(serviceNewCmd)
	root.AddCommand(serviceCmd)
}

func CleanService(project project.Project) {
	count := 0
	err := filepath.WalkDir(filepath.Join(project.BasePath, "api"), func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".pb.go") ||
			strings.HasSuffix(path, ".swagger.json") ||
			strings.HasSuffix(path, ".pb.validate.go") {
			err := os.Remove(path)
			if err == nil {
				count += 1
			}
			return err
		}
		return nil
	})

	if err == nil {
		// check generate file
		files := []string{"godepgraph.png", "mem.prof", "cpu.prof", "dep.dot", "openapi.yaml", ".log"}
		entries, readErr := os.ReadDir(".")

		if readErr == nil {
			for _, e := range entries {
				if e.IsDir() {
					continue
				}

				for _, f := range files {
					if strings.Index(e.Name(), f) == -1 {
						continue
					}
					_ = os.Remove(filepath.Join("./", e.Name()))
				}
			}
		} else {
			err = readErr
		}
	}

	if err != nil {
		s := fmt.Sprintf("cleaning project %s with error %s\n", project.ProjectName, err)
		panic(s)
	}

	fmt.Printf("cleaning project %s finished, %d files in directory %s has been cleaned\n", project.ProjectName, count, project.BasePath)
}

func InitialService(project project.Project) {
	cmdList := []*exec.Cmd{
		exec.Command("make", "proto-gen", fmt.Sprintf("service=%s", project.ServiceShortName)),
		exec.Command("goimports", "-w", "api"),
		exec.Command("go", "mod", "tidy"),
	}

	for _, cmd := range cmdList {
		cmd.Dir = project.BasePath
		err := cmd.Run()
		if err != nil {
			fmt.Printf("calling cmd %s with error %s\n", cmd.String(), err)
			os.Exit(0)
		}
	}
}

func CreateService() (project.Project, error) {
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

	duplicate := make([]string, 0)
	for _, s := range proj.ServiceList {
		duplicate = append(duplicate, s.Name)
		duplicate = append(duplicate, s.ShortName)
	}

	if err == nil {
		// service
		prompt = promptui.Prompt{
			Label:    "service",
			Default:  proj.ProjectName,
			Validate: Validator(duplicate, "service"),
		}
		p, err = prompt.Run()
		if strings.Index(p, "-service") == -1 {
			p = p + "-service"
		}
		proj.ServiceName = p
	}

	if err == nil {
		// module
		prompt = promptui.Prompt{
			Label:    "module",
			Default:  strings.ReplaceAll(proj.ProjectName, "-service", ""),
			Validate: Validator(nil, "module"),
		}
		p, err = prompt.Run()
		proj.ModuleName = p
	}

	proj.Prepare()
	return proj, err
}
