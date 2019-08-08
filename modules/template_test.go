package modules

import (
	"bufio"
	"os"
	"testing"

	"github.com/kireledan/gonkey/utils"
)

func TestArgParse(t *testing.T) {

	toread := "./test_assets/testargs.yaml"

	list := readArgs(toread)

	if list.Args["blah"] != "my arg" {
		t.Error("parsed incorrectly")
	}

	if list.Args["another_one"] != "23" {
		t.Error("number parsed wrong!")
	}

	if list.Args["mylastarg"] != "duh" {
		t.Error("parsed incorrectly")
	}

}

func TestTemplateRender(t *testing.T) {

	toread := "./test_assets/testtemplateargs.yaml"

	list := readArgs(toread)

	done := make(chan utils.Result)

	go renderfile("./test_assets/testtemplate.conf", "./test_assets/testresult.conf", list.Args, done)

	_ = <-done

	if _, err := os.Stat("./test_assets/testresult.conf"); os.IsExist(err) {
		t.Error("template not rendered")
	}

	file, err := os.Open("./test_assets/testresult.conf")
	if err != nil {

	}

	testfile := bufio.NewReader(file)

	// CHECK THE LINES CHANGED
	line, _ := Readln(testfile)
	if line != "export TODO_DIR=mydir" {
		t.Errorf("Line 1 not rendered")
	}
	line, _ = Readln(testfile)
	if line != "export TODO_FILE=anotheffile" {
		t.Errorf("Line 2 not rendered")
	}
	line, _ = Readln(testfile)
	if line != "export DONE_FILE=donefile" {
		t.Errorf("Line 3 not rendered")
	}
	line, _ = Readln(testfile)
	if line != "export REPORT_FILE=reportfile" {
		t.Errorf("Line 4 wasnt rendered.")
	}

	os.Remove("./test_assets/testresult.conf")

}
