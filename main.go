package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

type data struct {
	Dir string `json:"dir"`
	Command string `json:"cmd"`
	User string `json:"user"`
}

func shell(cmd string) ([]byte, error){
	c := exec.Command(sh, "-c", cmd)
	return c.CombinedOutput()
}

func hook(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		d := &data{}
		b,err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		err = json.Unmarshal(b, d)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		c := fmt.Sprintf("cd %s && %s", d.Dir, d.Command)
		out, _ := shell(c)

		w.Write(out)
		return
	} else if r.Method == http.MethodGet {
		dir := r.FormValue("dir")
		cmd := r.FormValue("cmd")
		user := r.FormValue("user")
		c := fmt.Sprintf("cd %s && sudo -u %s git %s", dir, user, cmd)
		out, _ := shell(c)

		w.Write(out)
		return
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

var addr string
var sh string

func init() {
	flag.StringVar(&addr, "l", ":10009", "listen address; like :10009")
	flag.StringVar(&sh, "s", "/bin/bash", "bash path")
}

func main() {
	flag.Parse()
	http.HandleFunc("/", hook)
	log.Fatal(http.ListenAndServe(addr, nil))
}