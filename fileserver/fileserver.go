package main

import (
	"net/http"
	"github.com/xkenmon/go-tools/fileserver/controller"
	"github.com/xkenmon/wego"
	"log"
	"github.com/xkenmon/wego/router"
	"github.com/xkenmon/wego/context"
)

func main() {
	app := wego.NewApp()
	app.HandleFunc("GET", "/list/**", controller.ListFile).
		Handle("Get", "/static/**", router.FileHandler(http.Dir(""))).
		HandleFunc("Get", "/download/**", controller.DownloadZip).
		HandleFunc("Get", "/*", controller.Index).
		AddHookFunc("Get", "/download/*", router.HookBefore, func(ctx *context.Context) {
		log.Println("Downloading: ", ctx.Request.URL)
	}).
		AddHookFunc("Get", "/**", router.HookBefore, func(ctx *context.Context) {
		log.Println("Access: ", ctx.Request.URL)
	})
	log.Fatal(app.Host("localhost").Port(8080).Run())
}
