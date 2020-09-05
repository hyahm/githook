package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hyahm/goconfig"
	"github.com/hyahm/golog"
	"github.com/hyahm/xmux"
)

func hook(w http.ResponseWriter, r *http.Request) {
	golog.Info(r.Method)
	token := r.Header.Get("X-Gitlab-Token")
	golog.Info(token)
	filename := xmux.Var(r)["filename"]
	b, err := ioutil.ReadFile(filepath.Join(goconfig.ReadString("server.jsondir"), filename))
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	jf := &JsonFile{}
	err = json.Unmarshal(b, &jf)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	out, err := jf.shell()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(out)
	return

}

type JsonFile struct {
	User    string `json:"user"`
	Command string `json:"command"`
	Dir     string `json:"dir"`
	Shell   string `json:"shell"`
}

func (jf *JsonFile) shell() ([]byte, error) {
	c := exec.Command(jf.Shell, "-c", fmt.Sprintf("cd %s && sudo -u %s %s", jf.Dir, jf.User, jf.Command))
	return c.CombinedOutput()
}

func main() {
	conf := "hook.ini"
	if len(os.Args) > 1 {
		conf = os.Args[1]
	}
	goconfig.InitConf(conf, goconfig.INI)
	_, err := os.Stat(goconfig.ReadString("server.jsondir"))
	if os.IsNotExist(err) {
		if err = os.MkdirAll(goconfig.ReadString("server.jsondir"), 0755); err != nil {
			log.Fatal(err)
		}
	}
	router := xmux.NewRouter()
	router.SetHeader("Access-Control-Allow-Origin", "*")
	router.SetHeader("Content-Type", "application/x-www-form-urlencoded,application/json; charset=UTF-8")
	router.SetHeader("Access-Control-Allow-Headers", "Content-Type,Access-Token,X-Token,smail,X-Gitlab-Token")
	router.Post("/post/{filename}", hook)
	router.Get("/get/{filename}", hook)
	golog.Info("listen on ", goconfig.ReadString("server.listen", ":10009"))
	err = router.Run(goconfig.ReadString("server.listen", ":10009"))
	if err != nil {
		log.Fatal(err)
	}
}
