# go-pc_home

* 查看当前环境
echo $GIN_MODE
* 配置环境
```
    export GIN_MODE=debug
    export GIN_MODE=test
    //正式
    export GIN_MODE=release
```
* 运行
go run *.go

* 自动创建struct

```
    db2struct --host 10.10.1.23 -d chelun_home -t home_news --package models --struct homeNews -p --user chelun --guregu --gorm
```