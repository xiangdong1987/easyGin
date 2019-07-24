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
