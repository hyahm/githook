package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/hyahm/goconfig"
	"github.com/hyahm/golog"
)

type JsonFile struct {
	User    string `json:"user"`
	Command string `json:"command"`
	After   string `json:"after"`
	Dir     string `json:"dir"`
	Shell   string `json:"shell"`
	Env     string `json:"env"`
}

func read(rc io.ReadCloser, iserr bool) {
	br := bufio.NewReader(rc)
	for {
		line, _, err := br.ReadLine()
		if err != nil {
			return
		}
		if iserr {
			golog.Error(string(line))
		} else {
			golog.Info(string(line))
		}

	}

}

func (jf *JsonFile) shell() error {
	err := jf.cmd(jf.Command)
	if err != nil {
		return err
	}
	go jf.cmd(jf.After)
	return nil
	// return c.CombinedOutput()
}

func (jf *JsonFile) cmd(cmd string) error {

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
		c = exec.Command(jf.Shell, arg, cmd)
	} else {
		c = exec.Command(jf.Shell, arg, fmt.Sprintf("sudo -u %s %s", jf.User, cmd))
	}
	c.Dir = jf.Dir
	sep, err := c.StderrPipe()
	if err != nil {
		golog.Error(err)
	}
	go read(sep, true)
	sop, err := c.StdoutPipe()
	if err != nil {
		golog.Error(err)
	}
	go read(sop, false)
	c.Env = os.Environ()
	golog.Info(c.Env)
	err = c.Start()
	if err != nil {
		golog.Error(err)
	}

	return c.Wait()
}

func pull(filename string) error {

	b, err := ioutil.ReadFile(filepath.Join(goconfig.ReadString("jsondir"), filename))
	if err != nil {
		return err
	}
	golog.Info(string(b))
	jf := &JsonFile{}
	err = json.Unmarshal(b, &jf)
	if err != nil {
		return err
	}
	err = jf.shell()
	if err != nil {
		golog.Error(err)
		return err
	}
	return nil
}
