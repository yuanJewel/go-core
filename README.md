# go-core

本库中封装了常用功能，方便快速开发。推荐go项目引用本项目，快速实现安全的crud，可以根据需求保存审计信息等。

已实现的特性：

- 数据库模块封装
  - 目前支持mysql，未来会支持其他数据库
  - 支持表不存在时自动创建表。
  - 支持数据记录，支持数据加密，支持数据校验。
  - 支持数据权限控制，支持数据审计。
  - 支持数据导出，支持数据导入。
  - 支持数据分页，支持数据排序，支持数据搜索。
  - 支持操作预期和事务控制，保障数据安全。
- api接口模块封装
  - 支持jwt认证，并可以刷新认证信息。
  - 支持数据传输标准化，包括输入输出的格式和数据校验。
  - 支持微服务架构的安全和鉴权。
  - 支持快速和数据库建立crud接口群。
- 日志模块封装
  - 统一所有日志输出方式，支持输出到文件，支持输出到控制台。
  - 支持日志文件自动切割，滚动迭代，按日期归档。
  - 数据库、api日志都可以定义和配置，包括debug级别的日志。
- 配置模块封装
  - 支持yaml和json两种格式，支持环境变量覆盖。

## 使用方法

### 直接运行Demo

本地运行该项目代码，需要先初始化，启动后，即可访问

```shell
make init   # 初始化本地环境，下载依赖包
make run    # 运行任务
```

### 作为库引用

可以仿照库中pkg目录使用该包，下载方法 `go get github.com/yuanJewel/go-core`

一个简单的web项目可以仿照如下：

```
├── asset               (go-bindata 自动生成)
├── docs                (swag init 自动生成)
├── Dockerfile          (构建镜像，可以参照本项目)
├── Makefile            (编译文件，可以参照本项目)
├── README.md
├── application.yml     (默认配置文件，不建议提交到git中)
├── go.mod              (go mod init 自动生成)
├── main.go             (入口文件，可以参照本项目)
├── pkg                 (集成项目，自定义自己web项目逻辑的位置，结构可以参照本项目)
│    ├── api            (接口类的实现)
│    ├── config         (配置类的实现)
│    └── db             (数据库类的实现)
└── views               (自定义需要打包进过项目的前端页面，目前只能应用于'/api/v1'下的路由使用)
    ├── 404.html
    └── index.html
```

### 环境变量

- `RECORD_DATA`: 是否开启数据记录，默认为true
- `CHECK_TABLE_EXISTS`: 是否检查表是否存在，(需要数据库高级权限)，默认为false
- `LOGGER_RETAIN_NUMBER`: 日志文件保留个数，默认为3
- `LOGGER_FILE_SIZE`: 日志文件最大大小，默认为50M
- `LOGGER_ROOT_PATH`: 日志文件存放位置，默认为`./logs`
- `LOGGER_OUT_STYLE`: 日志输出方式，默认为`file`，可以选择`stdout`不输出到文件
- `LOGGER_ACCESS_OUT_STYLE`: api访问日志输出格式，默认为`file`，可以选择`stdout`不输出到文件
- `LOGGER_OUT_LEVEL`: 日志输出级别，默认为`info`

### 配置基础格式

```yaml
apiVersion: v1

server:
  port: 8080

auth:
  key:                  # jwt的密钥
  timeout: 600
  refresh: 300
  cryptoKey:            # 数据加密密钥
  cryptoPrefix:         # 数据加密前缀

db:
  # 目前只支持mysql
  driver: mysql
  host:
  port: 3306
  db:
  user:
  password:
  charset: utf8
  idle_connections: 2
  max_connections: 10
```

## 目录结构

### api

目录 `api` 中声明，引用iris框架，实现api接口，封装了认证基本接口，安全检查接口

### db

目录 `db` 中声明包括数据库目录

- `object` 数据库的数据结构
- `service` 数据库接口
- `mysql` 数据库具体方法

### 日志

目录 `logger` 中声明，引用logrus库实现日志的基本功能

## 声明

引用或商用请注明出处，详细授权使用声明请查看: [LICENSE](https://github.com/yuanJewel/go-core/blob/main/LICENSE)

大版本间隔可能会存在不兼容情况，更新迭代，请自行评估。

如果有扩展需求或使用问题，请提交issue，或者邮件联系作者：`luyu151111@163.com`