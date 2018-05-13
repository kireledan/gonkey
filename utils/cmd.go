package utils

import (
	"bytes"
	"fmt"
	"github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func readKey(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func execRemoteCommandHost(cmd string, h Host) string {
	return execRemoteCommand(cmd, h.ip, h.user, h.key)
}

type Result struct {
	RC      int
	Stdout  string
	Stderr  string
	Command string
	Args    []string
}

func (m *Result) GetRC() int {
	return m.RC
}

func (m *Result) GetStdout() string {
	return m.Stdout
}

func (m *Result) GetStderr() string {
	return m.Stderr
}

const defaultFailedCode = 1

func ExecCmd(cmd string, done chan Result) {
	parts := strings.Fields(cmd)

	cmd_exec := exec.Command("bash", "-c", cmd)
	var outbuf, errbuf bytes.Buffer
	cmd_exec.Stdout = &outbuf
	cmd_exec.Stderr = &errbuf

	res := new(Result)

	err := cmd_exec.Run()
	res.Command = cmd
	res.Args = parts
	res.Stdout = outbuf.String()
	res.Stderr = errbuf.String()

	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			res.RC = ws.ExitStatus()
		} else {
			// This will happen (in OSX) if `name` is not available in $PATH,
			// in this situation, exit code could not be get, and stderr will be
			// empty string very likely, so we use the default fail code, and format err
			// to string and set to stderr
			log.Printf("Could not get exit code for failed program: %v", cmd)
			res.RC = defaultFailedCode
			if res.Stderr == "" {
				res.Stderr = err.Error()
			}
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd_exec.ProcessState.Sys().(syscall.WaitStatus)
		res.RC = ws.ExitStatus()
	}

	done <- *res
}

func execRemoteCommand(cmd string, server string, user string, key string) string {

	clientConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			readKey(key),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client := scp.NewClient(server, clientConfig)

	err := client.Connect()
	if err != nil {
		fmt.Println("Couldn't establish a connection to the remote server ", err)
		return "error"
	}

	ret, err := client.Session.CombinedOutput(cmd)
	s := string(ret)

	if err != nil {
		fmt.Println("ERR: ", s)
		return err.Error()
	}

	return s

}

func copyFileToServerHost(source string, dest string, permissions string, h Host) {
	copyFileToServer(source, dest, h.ip, h.user, h.key, permissions)
}

func copyFileToServer(source string, dest string, server string, user string, key string, permissions string) {

	clientConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			readKey(key),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client := scp.NewClient(server, clientConfig)

	err := client.Connect()
	if err != nil {
		fmt.Println("Couldn't establish a connection to the remote server ", err)
		return
	}

	f, _ := os.Open(source)
	defer client.Session.Close()
	defer f.Close()
	client.CopyFile(f, dest, permissions)

}

func GetBinaryPath(bin string) (string, bool) {

	pathSeparator := string(os.PathListSeparator)

	path := os.Getenv("PATH")
	systemPaths := append(strings.Split(path, pathSeparator), "/sbin", "/usr/sbin", "/usr/local/sbin")

	for _, element := range systemPaths {
		if _, err := os.Stat(element + "/" + bin); os.IsNotExist(err) {
			continue
		} else {
			return element + "/" + bin, true
		}
	}
	return "Not found", false
}
