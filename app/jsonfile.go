package app

import (
	"fmt"
	"os/exec"
	"runtime"
)

type JsonFile struct {
	User    string `json:"user"`
	Command string `json:"command"`
	Dir     string `json:"dir"`
	Shell   string `json:"shell"`
}

func (jf *JsonFile) shell() ([]byte, error) {
	var c *exec.Cmd
	arg := "-c"
	if jf.Shell == "" {
		jf.Shell = "/bin/bash"
		if runtime.GOOS == "windows" {
			arg = "/c"
			jf.Shell = "powershell"
			jf.User = ""
		}
	}

	if jf.User == "" {
		c = exec.Command(jf.Shell, arg, jf.Command)
	} else {
		c = exec.Command(jf.Shell, arg, fmt.Sprintf("sudo -u %s %s", jf.User, jf.Command))
	}

	c.Dir = jf.Dir

	return c.CombinedOutput()
}
