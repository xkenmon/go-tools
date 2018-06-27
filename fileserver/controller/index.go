package controller

import (
	"net/http"
	"github.com/xkenmon/wego"
)

func Index(w http.ResponseWriter, r *http.Request) {

	p := r.FormValue("contextPath")
	if p != "" {
		u := r.URL
		u.Path = "/list"
		u.RawQuery = ""

		SetCtxPath(p)
		http.Redirect(w, r, u.String(), http.StatusSeeOther)
		return
	}

	tpl, err := wego.GetTpl("index.tpl")
	checkErr(err)
	err = tpl.Execute(w, nil)
	checkErr(err)

	defer WriteErr(w, recover())
}
