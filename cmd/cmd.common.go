package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/manifoldco/promptui"

	"github.com/eden-quan/eden-cli/project"
)

func chooseService(args []string, msg string) (service string, err error) {
	proj, err := NewProject()
	if err != nil {
		return "", err
	}

	service = ""

	names := make([]string, 0)
	for _, svc := range proj.ServiceList {
		names = append(names, svc.Name)
	}

	if msg == "" {
		msg = "please select the service needed to run"
	}

	if len(args) == 0 {
		if len(names) == 1 {
			service = names[0]
		} else {
			// use option
			sel := promptui.Select{
				Label: msg,
				Items: names,
				Templates: &promptui.SelectTemplates{
					Active:   "[{{.}}]",
					Inactive: "{{.}}",
					Selected: "Selected: {{.}}",
				},
			}
			_, service, err = sel.Run()
			if err != nil {
				err = fmt.Errorf("show selection failed with error %s", err)
				return
			}
		}
	} else if len(args) == 1 {
		for _, svc := range names {
			if strings.Index(svc, args[0]) != -1 {
				service = svc
			}
		}
		if service == "" {
			err = fmt.Errorf("can not found the service named like %s", args[0])
			return
		}
	} else {
		err = fmt.Errorf("only one service can run at a time, but we got %v", args)
		return
	}

	if service == "" {
		err = errors.New("service name is empty")
	}

	return
}

func NewProject() (project.Project, error) {
	proj := project.Project{}
	basePath, _ := filepath.Abs(".")

	// check if we are in the project
	modFile, err := os.OpenFile(path.Join(basePath, "go.mod"), os.O_RDWR, 0644)
	if err != nil {
		return proj, errors.New("please step in the project directory first")
	}

	// we in the project, initialize project path and
	proj.BasePath = filepath.Dir(basePath)
	proj.ProjectName = filepath.Base(basePath)

	// read package name
	reader := bufio.NewReader(modFile)
	line, _, _ := reader.ReadLine()
	packageName := string(line)[7:]
	proj.PackageName = strings.TrimSpace(packageName)
	err = proj.ScanServices()

	return proj, nil
}

func Validator(duplicateList []string, name string) func(string) error {
	return func(s string) error {
		if len(s) == 0 {
			return fmt.Errorf("%s name is empty", name)
		}

		if strings.Index(s, " ") != -1 {
			return fmt.Errorf("%s name with space", name)
		}

		if duplicateList != nil && slices.Index(duplicateList, s) != -1 {
			return fmt.Errorf("%s already exists", name)
		}
		return nil
	}
}
