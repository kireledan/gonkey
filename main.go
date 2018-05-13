package main

import (
	"github.com/kireledan/gonkey/core"
	"os"
)

func main() {
	pipeline := core.CreateSerialPipeline(os.Args[1])

	go pipeline.Run(nil)
	err := pipeline.Wait()
	if err != nil {
		println("pipeline failed")
	} else {
		println("pipeline ran")
	}
}
