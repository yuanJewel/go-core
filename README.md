# go-core

## 使用方法

```shell
make init   # 初始化本地环境，下载依赖包
make run    # 运行任务
```

## 目录结构

### api

目录 `api` 中声明，引用iris框架，实现api接口，封装了认证基本接口，安全检查接口

### db

包括数据库目录
- `db` 数据库的数据结构
- `cmdb` 数据库接口
- `mysql` 数据库具体方法

### 日志

目录 `logger` 中声明，引用logrus库实现日志的基本功能