package main

import (
	"net/http"
	"github.com/xkenmon/go-tools/fileserver/controller"
	"github.com/xkenmon/wego"
	"log"
	"github.com/xkenmon/wego/router"
)

func main() {
	app := wego.NewApp()
	app.HandleFunc("GET", "^/list/*", controller.ListFile).
		Handle("Get", "^/static/*", http.FileServer(http.Dir(""))).
		HandleFunc("Get", "^/download/*", controller.DownloadZip).
		HandleFunc("Get", "\\w*", controller.Index).
		AddHookFunc("Get", "/download", router.HookBefore, func(w http.ResponseWriter, r *http.Request) {
		log.Println("Downloading: ", r.URL.Path)
	}).
		AddHookFunc("Get", "/*", router.HookBefore, func(w http.ResponseWriter, r *http.Request) {
		log.Println("Access: ", r.URL)
	})
	log.Fatal(app.Host("localhost").Port(8080).Run())
}
