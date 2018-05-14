package main

import (
	"github.com/kireledan/gonkey/core"
	"os"
	"fmt"
)

func main() {
	pipeline := core.CreateSerialPipeline(os.Args[1])

	if pipeline == nil {
		fmt.Println("The pipeline was unable to be parsed. Please see above.")
		return
	}

	go pipeline.Run(nil)
	err := pipeline.Wait()
	if err != nil {
		println("pipeline failed!")
	}
}
