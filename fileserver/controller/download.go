package controller

import (
	"github.com/xkenmon/wego/util"
	"path/filepath"
	"os"
	"github.com/xkenmon/wego/context"
	"errors"
)

func DownloadZip(ctx *context.Context) error {
	paths, err := util.ParsePathVar(ctx.Request, "/download/{**}")
	ctxPath := ctx.In.GetSessionOrElse("ctxPath", ".").(string)
	if err != nil {
		return err
	}
	path := filepath.Clean(paths[0])
	if path == "" {
		path = "."
	}
	ctx.ResponseWriter.Header().Set("Content-Type", "application/zip")
	absPath := filepath.Join(ctxPath, path)
	if stat, err := os.Stat(absPath); err != nil || !stat.IsDir() {
		return errors.New(stat.Name() + " 不是目录")
	}
	if err = util.ZipDir(absPath, ctx.ResponseWriter); err != nil {
		return err
	}
	return nil
}
