package core

import (
	"errors"
	"os"

	"github.com/imkira/go-task"
	"github.com/kireledan/gonkey/modules"
)

func createTask(m ssTask) task.Task {
	moduleastask := func(t task.Task, ctx task.Context) {
		println("RUNNING TASK =>>> " + m.Label)
		results := modules.RunModule(m.ModuleToRun)
		if results.GetRC() != 0 {
			os.Stderr.WriteString("TASK FAILED -> " + results.Stderr)
			t.Cancel(errors.New(results.GetStdout()))
		} else {
			println("STDOUT:::" + results.Stdout)
		}
		println("\n....................\n")
	}
	return task.NewTaskWithFunc(moduleastask)
}

// CreateSerialPipeline returns a serial group that can be executed
func CreateSerialPipeline(file string) *task.SerialGroup {

	taskgroup := readYaml(file)

	if taskgroup == nil {
		return nil
	}

	maintrack := task.NewSerialGroup()
	for _, t := range taskgroup {
		maintrack.AddChild(createTask(t))
	}

	return maintrack
}
