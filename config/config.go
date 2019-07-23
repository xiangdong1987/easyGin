package config

var MysqlConfigMap = map[string]map[string]string{
	"debug": {
		"test":  "chelun:chelun@tcp(10.10.1.23:3306)/chelun_home?charset=utf8&parseTime=True&loc=Local&timeout=10ms",
	},
}
