package utils

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"
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
	err     error
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

func ExecCmd(cmd string) (string, error) {
	cmd_exec := exec.Command("bash", "-c", cmd)
	var wg sync.WaitGroup
	var outbuf, errbuf bytes.Buffer

	stdout, err := cmd_exec.StdoutPipe()
	if err != nil {
		return
	}

	stderr, err := cmd_exec.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd_exec.Start(); err != nil {
		return err
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		copy(stdout, outbuf)
	}()

	go func() {
		defer wg.Done()
		copy(stderr, errbuf)
	}()

	wg.Wait()

	if err := cmd_exec.Wait(); err != nil {
		return err
	}

	return nil
}

func copy(r io.Reader, target bytes.Buffer) {
	buf := make([]byte, 80)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			target.Write(buf[0:n])
		}
		if err != nil {
			break
		}
	}
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
