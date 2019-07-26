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

```
go run main.go --ifScaffold 1 --struct structName --index db_index databaseName --table tableName
```
## 待开发（思考中，欢迎大神提供建议）

* 代码测试覆盖
* 性能测试

