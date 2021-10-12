package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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

func read(rc io.ReadCloser, iserr bool) {
	br := bufio.NewReader(rc)
	for {
		line, _, err := br.ReadLine()
		if err != nil {
			return
		}
		if iserr {
			golog.Error(line)
		} else {
			golog.Info(line)
		}

	}

}

func (jf *JsonFile) shell() error {
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

	err = c.Start()
	if err != nil {
		golog.Error(err)
	}

	return c.Wait()
	// return c.CombinedOutput()
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
