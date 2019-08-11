package core

import (
	"testing"

	"github.com/kireledan/gonkey/modules"
	"gotest.tools/assert"
)

func TestParsing(t *testing.T) {
	filename := "./test_assets/example.yaml"

	config := readTasks(filename)

	if config.Tasks[0].Name != "run a command" {
		t.Error("not parsed...")
	}

	if config.Tasks[1].Name != "another" {
		t.Error("incorrect")
	}

	if config.Tasks[1].Args["cmd"] == "touch blah2" {
		t.Error("incorrect map parse")
	}

}

func TestModuleParse(t *testing.T) {
	config := readTasks("./test_assets/example.yaml")
	parsedmodules, err := createSSTasks(config)
	assert.NilError(t, err, "error parsing task list")

	firstmodule := parsedmodules[0]

	if modules.GetModuleName(firstmodule.ModuleToRun) != "command" {
		t.Error("incorrect parsed module.")
	}

}

func TestBadModuleParse(t *testing.T) {
	config := readTasks("./test_assets/bad.yaml")
	_, err := createSSTasks(config)
	assert.ErrorContains(t, err, "Unable to parse yaml")
}
