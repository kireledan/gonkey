package modules

import (
	"io"
	"log"
	"os"
	"strconv"

	"github.com/kireledan/gonkey/utils"
)

type Copy struct {
	src   string
	dst   string
	mode  string
	owner string
	group string
}

func (m Copy) execute(args ...string) utils.Result {

	done := make(chan utils.Result)
	go copyfile(m.src, m.dst, m.mode, m.owner, m.group, done)

	result := <-done

	return result
}

func copyfile(src string, dst string, mode string, owner string, group string, done chan utils.Result) {

	res := new(utils.Result)

	originalFile, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer originalFile.Close()

	// Create new file
	newFile, err := os.Create(dst)
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()

	// Copy the bytes to destination from source
	bytesWritten, err := io.Copy(newFile, originalFile)
	if err != nil {
		log.Fatal(err)
	}

	// Commit the file contents
	// Flushes memory to disk
	err = newFile.Sync()
	if err != nil {
		log.Fatal(err)
	}

	utils.UpdateFilePerms(dst, mode, owner, group)

	res.Stdout = "Copied " + strconv.Itoa(int(bytesWritten)) + " bytes"

	done <- *res
}

func (m Copy) getModuleName() string {
	return "copy"
}

func (m Copy) InitModuleFromMap(args map[string]string) Module {
	m.src = args["src"]
	m.dst = args["dst"]
	m.mode = args["mode"]
	m.owner = args["owner"]
	m.group = args["group"]
	return m
}
