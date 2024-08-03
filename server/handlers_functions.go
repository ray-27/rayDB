// server/handlers_functions.go
package server

import (
    "embed"
    "html/template"
    "net/http"
    "path/filepath"
)

//go:embed views/*.html
var templatesFS embed.FS

func render_Template(w http.ResponseWriter, template_name string, data interface{}) {
    tmplPath := filepath.Join("views", template_name)
    tmpl, err := template.ParseFS(templatesFS, tmplPath)
    if err != nil {
        http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
        return
    }

    err = tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
    }
}
