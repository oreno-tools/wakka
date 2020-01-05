package main

import (
	"bytes"
	_ "fmt"
	"os/exec"
	"strings"
	"testing"
	_ "time"
)

func TestVersionFlag(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "-version")
	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout

	_ = cmd.Run()

	if !strings.Contains(stdout.String(), AppVersion) {
		t.Fatal("Failed Test")
	}
}
