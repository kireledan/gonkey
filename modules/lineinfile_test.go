package modules

import (
	"bufio"
	"bytes"
	"github.com/kireledan/gonkey/utils"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestFileReplacement(t *testing.T) {

	var buffer bytes.Buffer
	buffer.WriteString("NOCHANGE" + "\n")
	buffer.WriteString("NOCHANGE" + "\n")
	buffer.WriteString("NOCHANGE" + "\n")
	buffer.WriteString("CHANGEME in this line" + "\n")
	buffer.WriteString("NOCHANGE" + "\n")
	buffer.WriteString("NOCHANGE" + "\n")
	buffer.WriteString("NOCHANGE" + "\n")
	buffer.WriteString("NOCHANGE" + "\n")
	buffer.WriteString("NOCHANGE" + "\n")

	var err error
	err = ioutil.WriteFile("./test_assets/myfile", buffer.Bytes(), 0644)
	if err != nil {
		panic(err)
	}

	changes := changeset{}
	changes.Search = "CHANGEME"
	changes.Replace = "I was changed!"
	example := Lineinfile{}
	example.File = "./test_assets/myfile"
	example.Changes = changes

	channel := utils.Result{}

	channel = RunModule(example, "")

	if !strings.Contains(channel.GetStdout(), "./test_assets/myfile found and replaced match") {
		t.Errorf("change not found and replaced...")
	}

	file, err := os.Open("./test_assets/myfile")
	if err != nil {

	}
	testfile := bufio.NewReader(file)

	// CHECK THE LINES CHANGED
	line, _ := Readln(testfile)
	if !strings.Contains(line, "NOCHANGE") {
		t.Errorf("Line 1 not preserved")
	}
	line, _ = Readln(testfile)
	if !strings.Contains(line, "NOCHANGE") {
		t.Errorf("Line 2 not preserved")
	}
	line, _ = Readln(testfile)
	if !strings.Contains(line, "NOCHANGE") {
		t.Errorf("Line 3 not preserved")
	}
	line, _ = Readln(testfile)
	if line != "I was changed!" {
		t.Errorf("Line 4 wasnt changed.")
	}
	line, _ = Readln(testfile)
	if !strings.Contains(line, "NOCHANGE") {
		t.Errorf("Line 5 not preserved")
	}
	line, _ = Readln(testfile)
	if !strings.Contains(line, "NOCHANGE") {
		t.Errorf("Line 6 not preserved")
	}
	line, _ = Readln(testfile)
	if !strings.Contains(line, "NOCHANGE") {
		t.Errorf("Line 7 not preserved")
	}
	line, _ = Readln(testfile)
	if !strings.Contains(line, "NOCHANGE") {
		t.Errorf("Line 8 not preserved")
	}
	line, _ = Readln(testfile)
	if !strings.Contains(line, "NOCHANGE") {
		t.Errorf("Line 9 not preserved")
	}

	os.Remove("./test_assets/myfile")
}
