# Eden 脚手架

脚手架提供后台服务的基本管理功能，当前提供了项目的创建/微服务的查询/微服务创建/模块创建/代码生成/运行服务 等能力，通过脚手架
我们能够确保所有服务的结构保持一致，整合研发人员的经验。

## 安装

脚手架以 `go` 包的形式提供，可以通过如下方式安装

```shell
go install github.com/eden/eden-cli@latest
```

## 使用示例

### 功能列表

通过 `help` 命令可以查看当前脚手架提供的功能, 如通过以下命令

```shell
eden-cli help
````

可以得到脚手架的整体命令列表:

```shell
Usage:
  eden-cli [command]

Available Commands:
  all-in-one  merge all services in this project to an single application
  completion  Generate the autocompletion script for the specified shell
  gen         generate codes by eden plugins
  help        Help about any command
  module      manage service module
  project     project manager
  run         run the target service
  service     service manager
  version     version of this application

Flags:
  -h, --help   help for eden-cli

Use "eden-cli [command] --help" for more information about a command.
```

### 创建项目

通过 `project` 命令可以创建新的项目

```shell
eden-cli project new
```

该命令会以交互式的方式来创建新项目，创建项目时会附带 demo 示例，示例中包含基本的框架功能说明，创建完成后的目录结构大致如下：

```
├── api
│  ├── common
│  │  ├── v1                   // 当前已基本不使用，为兼容旧功能保留，勿更改
│  │  └── v3                   // 通用定义
│  │     ├── enums
│  │     ├── resources         // 该目录包含通用的返回结果，通过 flatten 的方式为所有接口提供统一的返回结果类型
│  │     └── services          // 该目录提供了微服务名配置能力，通过该能力为微服务提供服务发现 / 自动加载注册中心配置等能力
│  └── some-service            // 由工具创建的微服务，所有的接口 / 数据 / 错误 / 枚举信息都在该目录下对应的子目录进行定义
│     └── v1                   // 版本信息为创建微服务时指定
│        ├── enums             // 用于定义枚举类型的目录
│        ├── errors            // 用于定义该微服务的错误信息
│        ├── resources         // 用于定义该微服务的消息类型
│        └── services          // 用于定义该微服务的接口
├── app
│  └── some-service            // 这里为微服务用于实现实际业务逻辑的位置
│     ├── cmd                  // 该目录包含了微服务的入口 / 启动函数，!!!通常情况不修改该目录
│     │  ├── export
│     │  └── some-service
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

整个框架创建完成后会提供一个使用依赖注入框架 Injection 作为基础能力的开发框架，
Inject 提供了默认的依赖注入，主要包括了各个中间件的注入，需要使用以下提供的中间件时，
需要先注入 config 模块，通过 config 模块，为中间件提供配置中心的配置信息，通过该函数可以得到以下注入信息

> 当然，这些信息框架基本已默认配置好，唯一必须要做的是去统一的配置中心配置新服务的信息，具体的说明我们可以参考 Wiki 的开发指南

1. Logger 注入，提供了基础的日志功能，可通过 log.Logger 得到
2. LoggerManager 注入，提供了日志库管理能力，可通过 *LogManager 得到
3. Redis 注入，提供缓存的访问及管理能力，可通过 kit.Redis 得到
4. MongoDB 注入，提供了 MongoDB 数据库的访问能力，可通过 kit.MongoDB 得到
5. MySQL 注入，提供了 MySQL 数据库的访问能力，可通过 kit.MySQL 得到
6. Messaging 注入，提供了基于 RabbitMQ 的消息队列能力, 可通过 kt.MessageQueue 得到
7. Tracing 注入，提供了全局的链路跟踪能力，所有通过依赖注入的客户端都能够自动得到链路跟踪的能力


### 服务管理

以下命令可以为当前的服务创建新的模块，以及查询当前服务中包含哪些微服务

```shell
eden-cli service ls     // 列出当前项目的微服务列表
eden-cli service new    // 创建一个新的微服务/模块
```
### 添加新模块/服务
在 创建/拉取 项目后，可通过下列命令添加新的业务模块/微服务, 在创建的过程中会提示输入微服务名，当微服务名已存在时可用于添加新的模块，
当指定的微服务不存在时，则会同时创建微服务及新的模块

```shell
eden-cli service new     // 创建新服务
eden-cli module new      // 创建新模块
```

### 生成代码

我们的微服务的接口 / 错误 / 资源等信息，都统一通过 `protocol buffer` 的形式定义，
因此修改完定义信息后我们需要通过生成命令 `gen` 来生成最新的代码

```shell
eden-cli gen [service-name]
```

`service-name` 为可选项，当提供时脚手架会生成与其匹配的微服务代码, 当不提供任何微服务名时，会以交互式的方式选择需要生成代码的微服务

> 如提供了为服务名，则会以部分匹配的方式查找对应的服务，如存在 some-service，则只需要输入 some,
> 当存在多个服务都匹配服务名时，则对应的服务都会生成

### 启动微服务

开发完 / 或想测试新创建服务的示例时，可以通过启动服务的命令 `run` 来运行服务，如上文所述，如服务是新服务，则运行前需要先到配置中心配置运行的地址跟端口信息

> 后续我们会提供新服务的本地执行方式，减少新服务启动所需的非必要信息

```shell
eden-cli run [service-name]
```

同 `gen` 命令，这里的服务名也可以以交互式的方式选择微服务或运行指定的微服务
