package modules

import (
	"github.com/kireledan/gonkey/utils"
)

/*

Any module that can be run needs to have an execute, a module name, and a method to initialize given a set of args.
Any module that implements this interface can be run by gonkey.

*/
type Module interface {
	execute(args ...string) utils.Result
	getModuleName() string
	InitModuleFromMap(args map[string]string) Module
}

func RunModule(mod Module, args ...string) utils.Result {
	return mod.execute(args...)
}

func GetModuleName(mod Module) string {
	return mod.getModuleName()
}

/* STARTING TEMPLATE FOR ANY NEW MODULE. BE SURE TO ADD TO THE TYPE FACTORY!!!

package modules

import (
	"fmt"
	"github.com/kireledan/gonkey/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"text/template"
)

type <MODULENAME> struct {
	<ARGS>
}

func (m <MODULENAME>) execute(args ...string) utils.Result {

	done := make(chan utils.Result)
	go mymainfunc(args,done)

	result := <-done

	return result
}


func (m <MODULENAME>) getModuleName() string {
	return "template"
}

func (m <MODULENAME>) InitModuleFromMap(args map[string]string) Module {
	//assign variables here
	return m
}

*/
