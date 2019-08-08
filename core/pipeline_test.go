package core

import (
	"os"
	"testing"

	"github.com/imkira/go-task"
)

func TestPipeline(t *testing.T) {

	m := make(map[string]string)
	m["cmd"] = "touch ./test_assets/blah1"

	mytask := ssTask{}
	mytask.ModuleToRun = mapStringToModule("command", m)
	mytask.Label = "Create blah file"
	createdTask := createTask(mytask)

	m["cmd"] = "echo test >> ./test_assets/blah1"
	mytask = ssTask{}
	mytask.ModuleToRun = mapStringToModule("command", m)
	mytask.Label = "Write to blah file"
	dependentTask := createTask(mytask)

	m["cmd"] = "rm -f ./test_assets/blah1"
	mytask = ssTask{}
	mytask.ModuleToRun = mapStringToModule("command", m)
	mytask.Label = "Remove blah file"
	finaltask := createTask(mytask)

	/// Everything is a serial group + concurrent group for now.

	maintrack := task.NewSerialGroup()

	maintrack.AddChild(createdTask)
	maintrack.AddChild(dependentTask)

	go maintrack.Run(nil)

	err2 := maintrack.Wait()
	if err2 != nil {
		t.Error("error")
	}

	if _, err := os.Stat("./test_assets/blah1"); os.IsNotExist(err) {
		t.Error("Command not executed")
	}

	finaltask.Run(nil)
	finaltask.Wait()
	if _, err := os.Stat("./test_assets/blah1"); os.IsExist(err) {
		t.Error("Command not executed")
	}

}

func TestBadPipeline(t *testing.T) {

	m := make(map[string]string)
	m["cmd"] = "echo hello"

	mytask := ssTask{}
	mytask.ModuleToRun = mapStringToModule("command", m)
	mytask.Label = "Create blah file"
	createdTask := createTask(mytask)

	m["cmd"] = "invalidcommand"
	mytask = ssTask{}
	mytask.ModuleToRun = mapStringToModule("command", m)
	mytask.Label = "Running invalid command"
	dependentTask := createTask(mytask)

	m["cmd"] = "rm -f ./test_assets/blah1"
	mytask = ssTask{}
	mytask.ModuleToRun = mapStringToModule("command", m)
	mytask.Label = "Remove blah file"
	finaltask := createTask(mytask)

	/// Everything is a serial group + concurrent group for now.

	maintrack := task.NewSerialGroup()

	maintrack.AddChild(createdTask)
	maintrack.AddChild(dependentTask)
	maintrack.AddChild(finaltask)

	go maintrack.Run(nil)

	err2 := maintrack.Wait()
	if err2 == nil {
		t.Error("Pipeline didnt fail...")
	}

}

func TestParsedPipeline(t *testing.T) {
	// Parse the 2 commands and check if they ran
	maintrack := CreateSerialPipeline("./test_assets/example.yaml")

	go maintrack.Run(nil)
	err := maintrack.Wait()
	if err != nil {
		t.Error("oops...")
	}
}
