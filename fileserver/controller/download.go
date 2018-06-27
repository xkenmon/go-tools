package controller

import (
	"net/http"
	"github.com/xkenmon/wego/util"
	"path/filepath"
	"os"
)

func DownloadZip(w http.ResponseWriter, req *http.Request) {
	paths, err := util.ParsePathVar(req, "/download/{**}")
	if err != nil {
		WriteErr(w, err)
		return
	}
	path := filepath.Clean(paths[0])
	if path == "" {
		path = "."
	}
	w.Header().Set("Content-Type", "application/zip")
	absPath := filepath.Join(ctxPath, path)
	if stat, err := os.Stat(absPath); err != nil || !stat.IsDir() {
		w.Write([]byte("不是目录"))
		return
	}
	if err = util.ZipDir(absPath, w); err != nil {
		WriteErr(w, err)
	}
}
