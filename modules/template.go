package modules

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/kireledan/gonkey/utils"
	"gopkg.in/yaml.v2"
)

type Config map[string]string

func (p Config) GetConfig(config_value string) string {
	return p[config_value]
}

type Template struct {
	template  string
	output    string
	valuefile string
}

type ArgList struct {
	Args map[string]string `yaml:"args"`
}

func (m Template) execute(args ...string) utils.Result {

	done := make(chan utils.Result)

	templateargs := readArgs(m.valuefile)
	go renderfile(m.template, m.output, templateargs.Args, done)

	result := <-done

	return result
}

func readArgs(filename string) ArgList {

	var args ArgList
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(source, &args)
	if err != nil {
		panic(err)
	}

	return args
}

func renderfile(file string, output string, values map[string]string, done chan utils.Result) {
	res := new(utils.Result)
	f, err := os.Create(output)
	if err != nil {
		fmt.Println("create file: ", err)
		return
	}

	t, _ := template.ParseFiles(file)
	err = t.Execute(f, Config(values))
	if err != nil {
		fmt.Println(err)
	}

	f.Close()

	res.Stdout = "Rendered template " + file + " to file " + output

	done <- *res
}

func (m Template) getModuleName() string {
	return "template"
}

func (m Template) InitModuleFromMap(args map[string]string) Module {
	m.output = args["output"]
	m.template = args["template"]
	m.valuefile = args["valuefile"]
	return m
}
