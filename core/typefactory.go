package core

import (
	"reflect"
	"strings"

	"github.com/kireledan/gonkey/modules"
)

var typeRegistry = make(map[string]reflect.Type)

/*
This is where any type is added to the type registry
*/
func init() {
	typeRegistry[strings.ToLower(reflect.TypeOf(modules.Command{}).Name())] = reflect.TypeOf(modules.Command{})
	typeRegistry[strings.ToLower(reflect.TypeOf(modules.Lineinfile{}).Name())] = reflect.TypeOf(modules.Lineinfile{})
	typeRegistry[strings.ToLower(reflect.TypeOf(modules.Service{}).Name())] = reflect.TypeOf(modules.Service{})
	typeRegistry[strings.ToLower(reflect.TypeOf(modules.Template{}).Name())] = reflect.TypeOf(modules.Template{})
	typeRegistry[strings.ToLower(reflect.TypeOf(modules.Copy{}).Name())] = reflect.TypeOf(modules.Copy{})
	typeRegistry[strings.ToLower(reflect.TypeOf(modules.Package{}).Name())] = reflect.TypeOf(modules.Package{})
}
