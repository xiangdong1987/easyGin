package test

import (
	"easyGin/scaffold"
	"testing"
)

func TestInitModels(t *testing.T) {
	scaffold.InitModels("person", "Person", "models", true, true, true)
}
func TestGenerateCURD(t *testing.T) {
	println(scaffold.GenerateCURD("Person", "id"))
}
func TestInitRouter(t *testing.T) {
	println(scaffold.InitRouter("Person"))
}
func TestGenerateApi(t *testing.T) {
	println(scaffold.GenerateApi("Person"))
}
func TestInitApi(t *testing.T) {
	println(scaffold.InitApi("Person"))
}
