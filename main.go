package main

import (
	"githook/app"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hyahm/goconfig"
	"github.com/hyahm/golog"
	"github.com/hyahm/xmux"
)

func main() {
	defer golog.Sync()
	conf := "hook.ini"
	if len(os.Args) > 1 {
		conf = os.Args[1]
	}
	goconfig.InitConf(conf, goconfig.INI)
	_, err := os.Stat(goconfig.ReadString("jsondir"))
	if os.IsNotExist(err) {
		if err = os.MkdirAll(goconfig.ReadString("jsondir"), 0755); err != nil {
			log.Fatal(err)
		}
	}
	router := xmux.NewRouter()
	router.SetHeader("Access-Control-Allow-Origin", "*")
	router.SetHeader("Content-Type", "application/x-www-form-urlencoded,application/json; charset=UTF-8")

	router.Post("/gitlab/{filename}", app.GitLabHook).
		SetHeader("Access-Control-Allow-Headers", "Content-Type,Access-Token,X-Token,X-Gitlab-Token")

	router.Post("/github/{filename}", app.GitHubHook)
	router.Post("/gitee/{filename}", app.GiteeHook)

	golog.Info("listen on ", goconfig.ReadString("listen", ":10009"))
	golog.Info(goconfig.ReadDuration("readtimeout", time.Second*30))
	svc := &http.Server{
		Addr:        goconfig.ReadString("listen", ":10009"),
		Handler:     router,
		ReadTimeout: goconfig.ReadDuration("readtimeout", time.Second*30),
	}
	log.Fatal(svc.ListenAndServe())

}
