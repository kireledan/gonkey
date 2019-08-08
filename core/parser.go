package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/kireledan/gonkey/modules"
	"github.com/kireledan/gonkey/utils"
	"gopkg.in/yaml.v2"
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
	parsedmodules, err := createSSTasks(config)
	if err != nil {
		return nil
	}
	return parsedmodules
}

func readTasks(filename string) TaskList {

	var tasks TaskList
	source, err := ioutil.ReadFile(filename)
	err = yaml.UnmarshalStrict(source, &tasks)
	if err != nil {
		fmt.Println(err)
		fmt.Print("\n\n-----------------------\n")
	}

	return tasks
}

// Returns modules from a yaml task list
func createSSTasks(tasks TaskList) ([]ssTask, error) {

	var ss []ssTask

	for _, t := range tasks.Tasks {
		toinsert := new(ssTask)
		moduleToRun := mapStringToModule(t.ModuleName, t.Args)
		if moduleToRun == nil {
			return nil, errors.New("Unable to parse yaml, unable to map YAML to a module.")
		}
		toinsert.ModuleToRun = moduleToRun
		toinsert.Label = t.Name
		ss = append(ss, *toinsert)
	}

	return ss, nil
}

func mapStringToModule(name string, args map[string]string) modules.Module {
	if name == "" {
		return nil
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

	return nil
}
