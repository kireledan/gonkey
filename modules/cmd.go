package modules

import (
	"github.com/kireledan/gonkey/utils"
)

type Command struct {
	command, args string
}

func (m Command) execute(args ...string) utils.Result {

	done := make(chan utils.Result)

	go utils.ExecCmd(m.command, done)

	result := <-done

	return result
}

func (m Command) getModuleName() string {
	return "command"
}

func (m Command) InitModuleFromMap(args map[string]string) Module {
	m.command = args["cmd"]
	return m
}
