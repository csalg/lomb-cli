package utils

import (
	"embed"
	"net/http"
	"text/template"
)

func RenderTemplate(w http.ResponseWriter, r *http.Request, fs *embed.FS, templateName string, model interface{}) {
	tmpl, err := template.ParseFS(fs, templateName)
	if err != nil {
		panic("getting template from fs:" + err.Error())
	}
	err = tmpl.Execute(w, model)
	if err != nil {
		panic("rendering template:" + err.Error())
	}
}
