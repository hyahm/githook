package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/hyahm/goconfig"
	"github.com/hyahm/golog"
	"github.com/hyahm/xmux"
)

func GitHubHook(w http.ResponseWriter, r *http.Request) {
	golog.Info(r.Method)
	for k, v := range r.Header {
		golog.Infof("%s:%s\n", k, strings.Join(v, ","))
	}
	// token := r.Header.Get("X-Gitlab-Token")
	// if token != goconfig.ReadString("token.gitlab", "123456") {
	// 	w.WriteHeader(http.StatusNetworkAuthenticationRequired)
	// 	return
	// }
	filename := xmux.Var(r)["filename"]
	x, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	golog.Info(string(x))
	b, err := ioutil.ReadFile(filepath.Join(goconfig.ReadString("server.jsondir"), filename))
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	golog.Info(string(b))
	jf := &JsonFile{}
	err = json.Unmarshal(b, &jf)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	out, err := jf.shell()
	if err != nil {
		golog.Error(err)
		golog.Error(string(out))
		w.Write([]byte(err.Error()))
		return
	}
	golog.Info(string(out))
	w.Write(out)
	return

}
