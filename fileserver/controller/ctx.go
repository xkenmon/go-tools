package controller

import "path/filepath"

var ctxPath = ""

func SetCtxPath(path string) {
	ctxPath = filepath.Clean(path)
}
