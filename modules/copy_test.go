package modules

import (
	"os"
	"testing"

	"github.com/kireledan/gonkey/utils"
)

func TestCopy(t *testing.T) {

	done := make(chan utils.Result)
	go copyfile("./test_assets/filetocopy.txt", "./test_assets/imacopy.txt", "", "", "", done)

	_ = <-done

	if _, err := os.Stat("./test_assets/imacopy.txt"); os.IsExist(err) {
		t.Error("template not rendered")
	}

	os.Remove("./test_assets/imacopy.txt")
}

func TestCopyMode(t *testing.T) {

	done := make(chan utils.Result)
	go copyfile("./test_assets/filetocopy.txt", "./test_assets/imacopyrestricted.txt", "0557", "en186015", "", done)

	_ = <-done

	if _, err := os.Stat("./test_assets/imacopyrestricted.txt"); os.IsExist(err) {
		t.Error("template not rendered")
	}

	info, _ := os.Stat("./test_assets/imacopyrestricted.txt")
	if info.Mode() != os.FileMode(0557) {
		t.Error("Mode not changed correctly")
	}

	os.Remove("./test_assets/imacopyrestricted.txt")
}
