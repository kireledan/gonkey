package utils

import (
	"bufio"
	"fmt"
	"os"
	s "strings"
)

// Host Struct to store individual host data
type Host struct {
	user string
	ip   string
	key  string
}

var (
	home   = os.Getenv("HOME")
	user   = os.Getenv("USER")
	gopath = os.Getenv("GOPATH")
)

func parseInventory(inventorypath string) []Host {
	var hosts []Host
	f, err := os.Open(inventorypath)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		os.Exit(1)
	}
	r := bufio.NewReader(f)
	s, e := Readln(r)
	for e == nil {
		hosts = append(hosts, parseLine(s))
		s, e = Readln(r)
	}
	return hosts
}

func parseLine(l string) Host {
	if len(l) > 0 {
		dat := s.Split(l, ",")
		ho := Host{ip: dat[0], user: s.TrimSpace(dat[1]), key: s.TrimSpace(dat[2])}
		return ho
	}
	return Host{}
}

// Readln returns a single line (without the ending \n)
// from the input buffered reader.
// An error is returned iff there is an error with the
// buffered reader.
func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix = true
		err      error
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}
