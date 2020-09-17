package app

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/hyahm/goconfig"
	"github.com/hyahm/golog"
	"github.com/hyahm/xmux"
)

func GiteeHook(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		golog.Infof("%s:%s\n", k, strings.Join(v, ","))
	}

	x, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	s := hmac.New(sha1.New, []byte(goconfig.ReadString("token", "123456")))
	s.Write(x)
	token := fmt.Sprintf("%x", s.Sum(nil))
	golog.Info(token)
	if "sha1="+token != r.Header.Get("X-Hub-Signature") {
		w.WriteHeader(http.StatusNetworkAuthenticationRequired)
		return
	}
	filename := xmux.Var(r)["filename"]

	w.Write(pull(filename))
	return

}
