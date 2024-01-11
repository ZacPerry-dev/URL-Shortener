package utils

import "html/template"

func ParseTemplates(files ...string) *template.Template {
	tmpl := template.New("")

	for _, path := range files {
		template.Must(tmpl.ParseFiles(path))
	}

	return tmpl
}
