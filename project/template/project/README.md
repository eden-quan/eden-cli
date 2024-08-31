# {{.StructServiceName }}

TODO: 补充项目说明

## 安装Golang

> Version 1.21.3

下载地址 [下载Golang](https://go.dev/dl/)

* windows安装 https://dl.google.com/go/go1.21.3.windows-amd64.msi
* mac-intel安装 https://dl.google.com/go/go1.21.3.darwin-amd64.pkg
* mac-Mx安装 https://dl.google.com/go/go1.21.3.darwin-arm64.pkg
* linux安装 https://dl.google.com/go/go1.21.3.linux-amd64.tar.gz

## 安装依赖工具

Eden 相关的项目都将使用 [eden-cli](https://github.com/eden-quan/eden-cli) 进行管理，该工具将为后端的项目提供统一的目录结构管理 / 服务维护 / 代码生成 / 运行程序 等基础能力.
鉴于各种网络 / Git 配置 / Go 私有仓库等原因，在运行任何项目前，请先完成如下配置

```shell
git config --global url."ssh://git@some_repo/".insteadOf "https://some_gitlab/"
go env -w GOPRIVATE=some_gitlab
```

如果手头上已有项目，并且尚未安装 `eden-cli` ，可通过在项目中执行 `make init` 进行依赖项的安装


## `eden-cli` 简单示例

整体项目的结构为: **项目 --> 服务 --> 模块**

### 创建新项目

执行如下命令后，工具会一步一步引导操作者在指定的目录中创建新项目, 在创建新项目的过程中，同时会为使用者创建默认的微服务及模块，
具体可参照项目的目录结构，具体的目录结构说明可参照 `eden-cli` 项目

```shell
eden-cli project new
```

### 添加新模块/服务
在 创建/拉取 项目后，可通过下列命令添加新的业务模块/微服务, 在创建的过程中会提示输入微服务名，当微服务名已存在时可用于添加新的模块，
当指定的微服务不存在时，则会同时创建微服务及新的模块

```shell
eden-cli service new     // 创建新服务
eden-cli module new      // 创建新模块
```

### 生成代码
微服务中的所有接口及消息类型，都会通过 `Protocol buffer` 进行定义，因此 `eden-cli` 工具同时也提供了快捷生成代码的命令，

```shell
eden-cli gen [service-name]
```

其中 `service-name` 为需要生成代码的微服务名，该名称支持模糊匹配，
意味着在需要生成 `some-service` 的代码时，可通过 `eden-cli gen some` 来触发生成功能。
当 `service-name` 为空时，会列出当前项目中包含的微服务列表，由用户挑选需要生成的服务。

### 运行项目

`eden-cli` 提供了 `run` 命令来运行指定的微服务，该命令支持与 `gen` 相同的模糊匹配逻辑。

```shell
eden-cli run [service-name]
```


## 项目结构

以下为本项目的整体结构说明

```
├── api
│  ├── common
│  │  ├── v1                   // 当前已基本不使用，为兼容旧功能保留，勿更改
│  │  └── v3                   // 通用定义
│  │     ├── enums
│  │     ├── resources         // 该目录包含通用的返回结果，通过 flatten 的方式为所有接口提供统一的返回结果类型
│  │     └── services          // 该目录提供了微服务名配置能力，通过该能力为微服务提供服务发现 / 自动加载注册中心配置等能力
│  └── {{.ServiceName}}        // 由工具创建的微服务，所有的接口 / 数据 / 错误 / 枚举信息都在该目录下对应的子目录进行定义
│     └── v1                   // 版本信息为创建微服务时指定
│        ├── enums             // 用于定义枚举类型的目录
│        ├── errors            // 用于定义该微服务的错误信息
│        ├── resources         // 用于定义该微服务的消息类型
│        └── services          // 用于定义该微服务的接口
├── app
│  └── {{.ServiceName}}            // 这里为微服务用于实现实际业务逻辑的位置
│     ├── cmd                  // 该目录包含了微服务的入口 / 启动函数，!!!通常情况不修改该目录
│     │  ├── export
│     │  └── {{.ServiceName}} 
│     ├── configs              // 该目录包含启动微服务所需的配置，一般情况下一个微服务只需包含配置中心的地址
│     │  └── config_app.yaml
│     ├── internal             // 微服务的业务逻辑实现地址
│     │  ├── conf              // 由工具生成的配置入口，一般情况下无需更改
│     │  ├── domain            // 每个模块各个功能领域的实现, 具体的业务逻辑在此处实现
│     │  ├── impl              // 实现 proto 的接口定义，一般负责注册接口到依赖框架及做参数校验
│     │  └── inject            // 依赖注入入口，需要注入新的组件时会引用该模块获取 Injector 进行注入
│     └── README.md
├── dep.dot                    // 由框架的依赖注入模块生成的模块依赖图
├── devops                     // 发布 / 构建脚本，一般无需关注该目录
│  ├── docker
│  └── makefile 
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── third_party                // 第三方的 Protocol 扩展，一般无需关注该目录
   ├── errors
   ├── google
   └── ...
```

## 依赖库

由 Eden CLI 创建的项目主要依赖下列公共库，虽然日常开发中无需使用到这些内容，但可以通过了解下列公共库的实现来理解整个项目的运行机制

- [Go Biz Kit](https://github.com/eden-quan/go-biz-kit)
  - 提供各种中间件的支持，如 MySQL、RabbitMQ、Tracing、Logging 等
  - 提供了通用的错误处理机制
  - 提供基于依赖注入的运行时支持
  - Web 及 Grpc 实现暂时通过 Kratos 实现 _(但开发人员不应对此做任何假设，因为后续可能会随时进行替换)_
- [Go Kratos Pkg](https://github.com/eden-quan/go-kratos-pkg)
  - 封装历史包 - 后续替换
- Tools
  - [protoc-gen-go-http-fx](https://github.com/eden-quan/protoc-gen-go-http-fx)
    - 扩展 http 插件，提供更通用的服务注册 / Flatten 等能力
  - [protoc-gen-go-grpc-fx](https://github.com/eden-quan/protoc-gen-go-grpc-fx)
    - 扩展 grpc 插件，提供更通用的服务注册 / Flatten 等能力
  - [protoc-gen-go-errors-fx](https://github.com/eden-quan/protoc-gen-go-errors-fx)
    - 扩展 errors 插件，提供结合 `BusinessKit` 提供通用的错误处理能力
  - [protoc-gen-openapi-fx](https://github.com/eden-quan/protoc-gen-openapi-fx) / [grpc-gateway](https://github.com/eden-quan/grpc-gateway)
    - 扩展 openapi 插件，提供与其他组件匹配的文档生成能力
  - [protobuf](https://github.com/eden-quan/protobuf)
    - 对 proto-go 插件进行扩展, 提供 Flatten 等能力
