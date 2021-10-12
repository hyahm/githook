package app

import (
	"net/http"

	"github.com/hyahm/goconfig"
	"github.com/hyahm/xmux"
)

func GiteeHook(w http.ResponseWriter, r *http.Request) {
	if goconfig.ReadString("token", "123456") != r.Header.Get("X-Gitee-Token") {
		w.WriteHeader(http.StatusNetworkAuthenticationRequired)
		return
	}
	filename := xmux.Var(r)["filename"]
	err := pull(filename)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

}
