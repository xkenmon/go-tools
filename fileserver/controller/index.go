package controller

import (
	"net/http"
	"github.com/xkenmon/wego"
	"github.com/xkenmon/wego/context"
)

func Index(ctx *context.Context) error {
	r := ctx.Request
	w := ctx.ResponseWriter

	ctxPath := ctx.In.GetSessionOrElse("ctxPath", ".").(string)

	p := r.FormValue("contextPath")
	if p != "" {
		u := r.URL
		u.Path = "/list/"
		u.RawQuery = ""

		ctx.In.SetSession("ctxPath", p)
		http.Redirect(w, r, u.String(), http.StatusSeeOther)
		return nil
	}

	tpl, err := wego.GetTpl("index.tpl")
	if err != nil {
		return err
	}
	return tpl.Execute(w, map[string]interface{}{
		"path": ctxPath,
	})
}
