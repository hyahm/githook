package main

import (
	"os/exec"
	"testing"
)

func TestWindows(t *testing.T) {
	cmd := exec.Command("powershell", "/c", "ls")
	out, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(out))
}
