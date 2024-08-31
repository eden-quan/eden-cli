package project

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/manifoldco/promptui"
)

//go:embed template
var EmbedProjectFS embed.FS

const (
	OverwriteDefault = 0
	OverwriteAllYes  = 1
	OverwriteALlNo   = 2
)

type ServiceName struct {
	Name      string
	ShortName string
}

type Project struct {
	BasePath              string // 构建目录
	ProjectName           string // 项目名，some-service, 非 -service 结尾的会添加 -service 后缀 Version             string // 版本号，v1
	Version               string // 版本号
	StructVersion         string // 驼峰版本号, 如会将 v1 转为 V1
	PackageName           string // 完整报名，github.com/eden/{project-name}
	ServiceName           string // 服务名 如 ping，最终会创建 ping-service
	StructServiceName     string // 服务名，用于定义结构体，因此需要为驼峰型, 如 ping-service 会转换为 PingService
	ServiceShortName      string // 服务名，缩减了 -service 部分, ping-service -> Ping
	ModuleName            string // 模块名，创建一个服务后，一个服务下可能会包含多个模块, 如 ping 服务下默认会有个 ping 模块
	StructModuleName      string // 驼峰模块名，如 ping 转为 Ping
	StructModuleUpperName string // 大写模块名
	EmbedPath             string
	FS                    *embed.FS
	BusinessPackageName   string // business 包名 github.com/eden/go-biz-kit
	OverwriteAll          int    // 文件重复时是否直接覆盖 0 为需要提问，1 为全部覆盖， 2 为全部不覆盖
	Empty                 string // 空占位符
	PlaceHolderIndex      int
	ServiceList           []ServiceName // 当前项目下的所有服务列表
}

// Exists 检查项目是否已存在
func (p *Project) Exists() bool {
	_, err := os.OpenFile(path.Join(p.BasePath, "go.mod"), os.O_RDWR, 0644)
	return err == nil
}

func (p *Project) CheckProjectName() {
	// TODO: 检查目录结构，满足条件时将当前目录名最为默认值
	lastIndex := strings.LastIndex(p.ProjectName, "-service")
	if lastIndex == -1 || p.ProjectName[lastIndex:] != "-service" {
		p.ProjectName += "-service"
	}
}

func (p *Project) Prepare() {

	// TODO: 将依赖库改为可配置，或是用户输入
	if len(p.BusinessPackageName) == 0 {
		p.BusinessPackageName = "github.com/eden/go-biz-kit"
	}

	p.ServiceShortName = strings.ReplaceAll(p.ServiceName, "-service", "")
	p.ServiceShortName = strings.ReplaceAll(p.ServiceShortName, "-", "")
	p.StructVersion = strcase.ToCamel(p.Version)
	p.StructServiceName = strcase.ToCamel(p.ServiceName)
	p.StructModuleName = strcase.ToCamel(p.ModuleName)
	p.StructModuleUpperName = strings.ToUpper(p.ModuleName)
	p.BasePath = path.Join(p.BasePath, p.ProjectName)

	content, err := os.ReadFile(path.Join(p.BasePath, ".eden-cli"))
	if err == nil {
		index, err := strconv.Atoi(string(content))
		if err == nil {
			p.PlaceHolderIndex = index
		}
	} else {
		p.PlaceHolderIndex = 99900
	}
}

func (p *Project) Done() {
	content := strconv.Itoa(p.PlaceHolderIndex)
	_ = os.WriteFile(path.Join(p.BasePath, ".eden-cli"), []byte(content), 0644)
}

func (p *Project) New(embedBasePath string) {

	p.FS = &EmbedProjectFS
	p.EmbedPath = embedBasePath

	Mkdir(p.BasePath)

	entries, err := p.FS.ReadDir(embedBasePath)
	if err != nil {
		panic(fmt.Sprintf("open template failed with error %s", err))
	}

	for _, e := range entries {
		p.Build(e, embedBasePath)
	}
}

func (p *Project) Build(d fs.DirEntry, basePath string) {
	if d.Name()[0] == '.' {
		fmt.Printf(d.Name())
	}
	if d.Type().IsRegular() { // file
		name := d.Name()
		filePath := path.Join(basePath, name)
		data, _ := p.FS.ReadFile(filePath)
		content := string(data)

		name, content = p.RunTemplate(filePath, content)
		p.Save(name, content)
	} else if d.Type().IsDir() {
		dirPath := path.Join(basePath, d.Name())
		entries, _ := p.FS.ReadDir(dirPath)
		for _, e := range entries {
			p.Build(e, dirPath)
		}
	}
}

func (p *Project) RunTemplate(name string, content string) (string, string) {
	// change name and content by template
	t := template.New(name)
	funcMap := template.FuncMap{
		// The name "title" is what the function will be called in the template text.
		"ServiceNameIndex": func() int {
			p.PlaceHolderIndex += 1
			return p.PlaceHolderIndex
		},
	}
	t.Funcs(funcMap)
	t, err := t.Parse(content)

	if err != nil {
		panic(fmt.Sprintf("Parse template %s with error %s\n", name, err))
	}

	contentBuf := bytes.NewBuffer([]byte{})
	err = t.Execute(contentBuf, p)
	if err != nil {
		fmt.Printf("parse file %s with error %s\n", name, err)
	}

	t = template.New(name)
	t, _ = t.Parse(name)
	nameBuf := bytes.NewBuffer([]byte{})
	err = t.Execute(nameBuf, p)
	if err != nil {
		fmt.Printf("parse file %s with error %s\n", name, err)
	}

	return nameBuf.String(), contentBuf.String()
}

func (p *Project) HasService(name string) bool {
	if len(p.ServiceList) == 0 {
		_ = p.ScanServices()
	}

	for _, s := range p.ServiceList {
		if s.Name == name || s.ShortName == name {
			return true
		}
	}

	return false
}

func (p *Project) ScanServices() error {
	services, err := os.ReadDir("./app")
	if err != nil {
		return fmt.Errorf("scan service failed with error %s\n", err)
	}

	names := make([]ServiceName, 0)
	for _, svc := range services {
		if strings.Index(svc.Name(), "-service") >= 0 {
			names = append(names, ServiceName{
				Name:      svc.Name(),
				ShortName: strings.ReplaceAll(svc.Name(), "-service", ""),
			})
		}
	}

	p.ServiceList = names
	return nil
}

func (p *Project) ScanModules(svc string) ([]string, error) {
	modules, err := os.ReadDir(fmt.Sprintf("./app/%s/internal/domain/repo", svc))
	if err != nil {
		return nil, fmt.Errorf("scan service failed with error %s\n", err)
	}

	names := make([]string, 0)
	for _, svc := range modules {
		if strings.Index(svc.Name(), ".repo.go") >= 0 {
			names = append(names, strings.ReplaceAll(svc.Name(), ".repo.go", ""))
		}
	}

	return names, nil
}

// Save 保存文件信息, 根据交互及配置决定是否覆盖
func (p *Project) Save(name string, content string) {
	flag := os.O_WRONLY | os.O_CREATE
	if p.OverwriteAll == OverwriteDefault || p.OverwriteAll == OverwriteALlNo {
		flag |= os.O_EXCL
	}
	name = name[len(p.EmbedPath):]
	target := path.Join(p.BasePath, name)
	baseDir := filepath.Dir(target)
	Mkdir(baseDir)
	f, err := os.OpenFile(target, flag, 0644)

	overwrite := true
	if os.IsExist(err) {
		if p.OverwriteAll == OverwriteDefault {
			prompt := promptui.Select{
				Label: fmt.Sprintf("file %s already exists, overwrite?", name),
				Items: []string{"Y", "N", "ALL(Y)", "ALL(N)"},
			}

			n, _, err := prompt.Run()
			if err != nil {
				os.Exit(0)
			}

			if n == 0 || n == 2 {
				overwrite = true
				if n == 2 {
					p.OverwriteAll = OverwriteAllYes
				}
			} else if n == 1 || n == 3 {
				overwrite = false
				if n == 3 {
					p.OverwriteAll = OverwriteALlNo
				}
			}
		} else if p.OverwriteAll == OverwriteAllYes {
			overwrite = true
		} else if p.OverwriteAll == OverwriteALlNo {
			overwrite = false
		}
	}

	if overwrite {
		flag = os.O_WRONLY | os.O_CREATE
		f, err = os.OpenFile(target, flag, 0644)
		if err == nil {
			err = f.Truncate(0)
			if err == nil {
				_, err = f.WriteString(content)
			}
		}
	} else {
		err = nil
	}

	if err != nil {
		fmt.Printf("create file %s failed with error %s\n", target, err)
		return
	}
}

func Mkdir(path string) {
	err := os.MkdirAll(path, 0777)
	if err != nil {
		fmt.Printf("create dir %s with error %s\n", path, err)
	}
}
