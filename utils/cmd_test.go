package utils

import (
	"strings"
	"testing"
)

func TestLocalExec(t *testing.T) {
	done := make(chan Result)

	/////////////////////  - Test a successful command
	go ExecCmd("echo blah", done)

	msg := <-done

	if strings.TrimRight(msg.Stdout, "\n") != ("blah") {
		t.Errorf("echo test failed")
	}

	if msg.RC != 0 {
		t.Errorf("Wrong error code, should be 0")
	}

	/////////////////////  - Test commands that dont exist

	go ExecCmd("lss ", done)

	msg = <-done
	if msg.RC != 127 {
		t.Errorf("Wrong error code returned")
	}
	if !strings.Contains(msg.Stderr, "lss: command not found") {
		t.Errorf("Error")
	}

	///////////////////// - Test invalid use of commands

	go ExecCmd("ls -l4", done)
	msg = <-done
	if strings.Split(msg.Stderr, "\n")[0] != "ls: illegal option -- 4" {
		t.Errorf("Error")
	}
	if msg.RC != 1 {
		t.Errorf("Wrong error code returned")
	}

}

func TestHostCommand(t *testing.T) {
	hosts := parseInventory("./test_assets/test_inventory")
	copyFileToServerHost("./test_assets/test", "/home/ubuntu/test", "0655", hosts[0])

	foundFiles := execRemoteCommandHost("ls /home/ubuntu/", hosts[0])
	if !strings.Contains(foundFiles, "test") {
		t.Errorf("test file not found in /home/ubuntu on host")
	}

	content := execRemoteCommandHost("cat /home/ubuntu/test", hosts[0])
	if !strings.Contains(content, "123") {
		t.Errorf("content of test file not copied to server")
	}
}

func TestBinFind(t *testing.T) {
	bin, result := GetBinaryPath("ssh")
	print(bin)
	if !result {
		t.Errorf("ssh bin not found")
	}
	if bin != "/usr/bin/ssh" {
		t.Errorf("ssh bin not found")
	}
}
