package modules

import (
	"github.com/kireledan/gonkey/utils"
)

type Service struct {
	name, status string
}

type Tools struct {
	serviceCommand, enableCommand string
}

func (m Service) execute(args ...string) utils.Result {
	finalCmd := generateRunCommand(m)

	done := make(chan utils.Result)

	go utils.ExecCmd(finalCmd, done)

	result := <-done

	return result
}

func generateRunCommand(m Service) string {
	serviceTools := getServiceTools()
	return serviceTools.serviceCommand + " " + serviceTools.enableCommand + " " + m.status + " " + m.name
}

// Get the tools to manage a service
// For now, its only systemd
func getServiceTools() Tools {

	servicebins := []string{"brew", "systemctl"}
	location := map[string]string{}
	availability := map[string]bool{}

	for _, element := range servicebins {
		location[element], availability[element] = utils.GetBinaryPath(element)
	}

	serviceCommand := "null"
	enableCommand := "null"

	if availability["systemctl"] {
		serviceCommand = location["systemctl"]
		enableCommand = ""
	}

	if availability["brew"] {
		serviceCommand = location["brew"]
		enableCommand = "services"
	}

	return Tools{serviceCommand, enableCommand}
}

func (m Service) getModuleName() string {
	return "service"
}

func (m Service) InitModuleFromMap(args map[string]string) Module {
	m.name = args["name"]
	m.status = args["status"]
	return m
}
