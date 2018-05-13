package modules

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/kireledan/gonkey/utils"
	"io/ioutil"
	"os"
	"regexp"
)

type changeset struct {
	SearchPlain string
	Search      string
	Replace     string
	MatchFound  bool
}

type changeresult struct {
	File   fileitem
	Output string
	Status bool
	Error  error
}

type fileitem struct {
	Path   string
	Output string
}

type Lineinfile struct {
	File    string
	Search  string
	Changes changeset
}

func (f Lineinfile) execute(args ...string) utils.Result {

	done := make(chan utils.Result)

	go runFileChange(f, done)

	result := <-done

	return result
}

func runFileChange(change Lineinfile, done chan utils.Result) {

	res := new(utils.Result)

	// Open file
	file, err := os.Open(change.File)
	if err != nil {

	}
	writeBufferToFile := false
	var buffer bytes.Buffer

	// Write file to buffer
	r := bufio.NewReader(file)
	line, e := Readln(r)
	for e == nil {
		newLine, lineChanged, skipLine := replaceLineInFile(line, change.Changes)

		if lineChanged || skipLine {
			writeBufferToFile = true
		}

		if !skipLine {
			buffer.WriteString(newLine + "\n")
		}
		line, e = Readln(r)
	}
	file.Close()

	var output string

	if writeBufferToFile {
		output = writeContentToFile(change.File, buffer)
		res.RC = 0
	} else {
		output = fmt.Sprintf("File %s contained no matches, not making changes", change.File)
		res.RC = 0
	}

	res.Stdout = output

	done <- *res
}
func replaceLineInFile(line string, changes changeset) (string, bool, bool) {
	changed := false
	skipLine := false

	if !searchMatch(line, changes) {
		changed = false
	} else {
		_ = regexp.MustCompile(changes.Search)
		line = changes.Replace

		changed = true
	}

	return line, changed, skipLine

}

// Checks if there is a match in content, based on search options
func searchMatch(content string, changeset changeset) bool {
	re := regexp.MustCompile(changeset.Search)
	if re.MatchString(content) {
		return true
	}

	return false
}

func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

func writeContentToFile(fileitem string, content bytes.Buffer) string {
	var err error
	err = ioutil.WriteFile(fileitem, content.Bytes(), 0644)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s found and replaced match\n", fileitem)

}

func (f Lineinfile) getModuleName() string {
	return "lineinfile"
}

func (f Lineinfile) InitModuleFromMap(args map[string]string) Module {
	changes := changeset{}
	changes.Search = args["search"]
	changes.Replace = args["replacewith"]
	f.File = args["file"]
	f.Changes = changes
	return f
}
