# easyGin
## 概述
* 实现了快速上手Gin，自动生成Restful Api的脚手架，代码依赖少轻量。可以快速开发Restful Api,代码一键生成。
## 依赖
* db2struck 代码生成对db2struck 进行了改造
* Gin 框架
* gorm mysql驱动
* go-redis redis 驱动
* go vendor 依赖管理
## 功能
* 项目的目录结构
* 快速restful api 生成
* go mod 接入
* go template 代码模板
* 代码测试覆盖
* 性能测试

## 使用go mod 自动加载
```
export GO111MODULE="on"
```


```
go run main.go --ifScaffold 1 --struct Person --index test --table person
```
## 测试
* 单元测试：对脚手架工具进行了单元测试
* 性能测试：虚拟机 2G 内存 单核 cup 结果如下
 ![设置1](/static/ab.png)

# docker mysql
```azure
docker run -itd --name mysql-test -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root mysql
```

# 自动生成open-api服务端
```azure
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=types.cfg.yaml person.ymal
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=server.cfg.yaml person.ymal
```