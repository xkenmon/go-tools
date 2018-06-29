package controller

import (
	"net/http"
	"os"
	"html/template"
	"github.com/dustin/go-humanize"
	"log"
	"github.com/xkenmon/wego/util"
	"github.com/xkenmon/wego"
	"path/filepath"
	"github.com/xkenmon/wego/context"
)

var funcList = template.FuncMap{
	"IsDir":   func(info os.FileInfo) bool { return info.IsDir() },
	"GetName": func(info os.FileInfo) string { return info.Name() },
	"GetTime": func(info os.FileInfo) string { return humanize.Time(info.ModTime()) },
	"GetSize": func(info os.FileInfo) string { return humanize.Bytes(uint64(info.Size())) },
}

func ListFile(ctx *context.Context) error {
	r := ctx.Request
	w := ctx.ResponseWriter

	ctxPath := ctx.In.GetSessionOrElse("ctxPath",".").(string)

	paths, err := util.ParsePathVar(r, "/list/{**}")
	if err != nil {
		return err
	}
	path := paths[0]
	log.Println(path)
	if path == "" {
		path = "."
	}
	dir, err := os.Open(filepath.Join(ctxPath, path))
	if info, err := dir.Stat(); err == nil {
		if !info.IsDir() {
			r.URL.Path = path
			http.FileServer(http.Dir(ctxPath)).ServeHTTP(w, r)
			return nil
		}
	} else {
		return err
	}
	infoList, err := dir.Readdir(-1)
	if err != nil {
		return nil
	}
	t, err := wego.GetTplWithFuncs(funcList, "list.tpl")
	if err != nil {
		return err
	}
	return t.Execute(w,
		map[string]interface{}{"info": infoList, "path": path, "parent": filepath.Join(path, "..")})
}
