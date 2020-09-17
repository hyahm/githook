package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/hyahm/goconfig"
	"github.com/hyahm/golog"
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

func pull(filename string) []byte {

	b, err := ioutil.ReadFile(filepath.Join(goconfig.ReadString("jsondir"), filename))
	if err != nil {
		return []byte(err.Error())
	}
	golog.Info(string(b))
	jf := &JsonFile{}
	err = json.Unmarshal(b, &jf)
	if err != nil {
		return []byte(err.Error())
	}
	out, err := jf.shell()
	if err != nil {
		golog.Error(err)
		return []byte(err.Error())
	}
	return out
}
