package modules

import (
	"github.com/kireledan/gonkey/utils"
	"runtime"
	"strings"
	"testing"
)

func TestServiceCommand(t *testing.T) {
	example := Service{}
	example.name = "myservice"
	example.status = "start"

	if runtime.GOOS == "darwin" {
		if "/usr/local/bin/brew services start myservice" != generateRunCommand(example) {
			t.Errorf("Wrong command")
		}
	} else {
		if "/bin/systemctl start myservice" != generateRunCommand(example) {
			t.Errorf("Wrong command")
		}
	}

	println(generateRunCommand(example))

}

func TestServiceRestart(t *testing.T) {
	example := Service{}
	example.name = "emacs"
	example.status = "restart"

	if runtime.GOOS == "darwin" {
		channel := utils.Result{}
		channel = RunModule(example, "")
		print(channel.GetStdout())
		if channel.GetRC() != 0 {
			t.Errorf("oops.")
		}

		if !strings.Contains(channel.GetStdout(), "==> Successfully started `emacs` (label: homebrew.mxcl.emacs)") {
			t.Errorf("Service not started")
		}
	}

}
