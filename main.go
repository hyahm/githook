package main

import (
	"flag"
	"fmt"
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
	w.Header().Set(header,token)
	if r.Method == http.MethodPost {
		c := fmt.Sprintf("cd %s && sudo -u %s git %s", dir, user, cmd)
		out, _ := shell(c)
		w.Write(out)
		return
	} else if r.Method == http.MethodGet {
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
var header string
var token string
var user string
var cmd string
var dir string

func init() {
	flag.StringVar(&addr, "l", ":10009", "listen address; like :10009")
	flag.StringVar(&sh, "s", "/bin/bash", "bash path")
	flag.StringVar(&header, "h", "X-Gitlab-Token", "header")
	flag.StringVar(&token, "t", "123456", "token")
	flag.StringVar(&cmd, "c", "pull", "cmd pull")
	flag.StringVar(&user, "u", "root", "cmd pull")
	flag.StringVar(&dir, "d", "/var/www/test", "dir name")
}

func main() {
	flag.Parse()
	http.HandleFunc("/", hook)
	log.Fatal(http.ListenAndServe(addr, nil))
}