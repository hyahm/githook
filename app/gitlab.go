package app

import (
	"net/http"

	"github.com/hyahm/goconfig"
	"github.com/hyahm/xmux"
)

func GitLabHook(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-Gitlab-Token")
	if token != goconfig.ReadString("token", "123456") {
		w.WriteHeader(http.StatusNetworkAuthenticationRequired)
		return
	}
	filename := xmux.Var(r)["filename"]
	err := pull(filename)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}
