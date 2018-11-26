package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun_alertFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("./pdump -alert", " ")

	status := cli.Run(args)
	_ = status
}

func TestRun_bufferFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("./pdump -buffer", " ")

	status := cli.Run(args)
	_ = status
}

func TestRun_secFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("./pdump -sec", " ")

	status := cli.Run(args)
	_ = status
}

func TestRun_nicFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("./pdump -nic", " ")

	status := cli.Run(args)
	_ = status
}

func TestRun_execFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("./pdump -exec", " ")

	status := cli.Run(args)
	_ = status
}
