package core

import (
	"github.com/kireledan/gonkey/modules"
	"github.com/kireledan/gonkey/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"fmt"
	"reflect"
)

type TaskList struct {
	Tasks []struct {
		Name       string            `yaml:"name"`
		ModuleName string            `yaml:"module"`
		Args       map[string]string `yaml:"args"`
		Register   string            `yaml:"register"`
	} `yaml:"tasks"`
}

type ssTask struct {
	ModuleToRun modules.Module
	Results     utils.Result
	ID          int
	Label       string
}

func readYaml(filename string) []ssTask {
	config := readTasks(filename)
	parsedmodules := createSSTasks(config)
	return parsedmodules
}

func readTasks(filename string) TaskList {

	var tasks TaskList
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		print("here?")
	}
	err = yaml.UnmarshalStrict(source, &tasks)
	if err != nil {
		fmt.Print(err)
	}

	return tasks
}

// Returns modules from a yaml task list
func createSSTasks(tasks TaskList) []ssTask {

	var ss []ssTask

	for _, t := range tasks.Tasks {
		toinsert := new(ssTask)
		toinsert.ModuleToRun = mapStringToModule(t.ModuleName, t.Args)
		toinsert.Label = t.Name
		ss = append(ss, *toinsert)
	}

	return ss
}

func mapStringToModule(name string, args map[string]string) modules.Module {
	if name == "" {
		panic("Task yaml incorrect.")
	}

	if val, ok := typeRegistry[name]; ok {
		if !ok {
			return nil
		}
		v := reflect.New(val).Elem()
		mod := v.Interface().(modules.Module)
		mod = mod.InitModuleFromMap(args)
		return mod
	}

	panic("INVALID MODULE " + name)
	return nil
}
