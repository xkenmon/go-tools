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
)

var funcList = template.FuncMap{
	"IsDir":   func(info os.FileInfo) bool { return info.IsDir() },
	"GetName": func(info os.FileInfo) string { return info.Name() },
	"GetTime": func(info os.FileInfo) string { return humanize.Time(info.ModTime()) },
	"GetSize": func(info os.FileInfo) string { return humanize.Bytes(uint64(info.Size())) },
}

func ListFile(w http.ResponseWriter, r *http.Request) {
	defer WriteErr(w, recover())
	paths, err := util.ParsePathVar(r, "/list/{**}")
	checkErr(err)
	path := paths[0]
	log.Println(path)
	if path == "" {
		path = "."
	}
	dir, err := os.Open(filepath.Join(ctxPath, path))
	if info, err := dir.Stat(); err == nil {
		if !info.IsDir() {
			r.URL.Path = path
			http.FileServer(http.Dir(".")).ServeHTTP(w, r)
			return
		}
	}
	checkErr(err)
	infoList, err := dir.Readdir(-1)
	checkErr(err)
	t, err := wego.GetTplWithFuncs(funcList, "list.tpl")
	err = t.Execute(w,
		map[string]interface{}{"info": infoList, "path": path, "parent": filepath.Join(path, "..")})
	checkErr(err)
}
