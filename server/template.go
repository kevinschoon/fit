package server

import (
	"github.com/kevinschoon/gofit/models"
	"html/template"
	"net/http"
)

const (
	staticDir string = "www"
	baseTmpl  string = staticDir + "/base.html"
	chartTmpl string = staticDir + "/chart.html"
	dataTmpl  string = staticDir + "/data.html"
)

type TemplateData struct {
	Columns    []string
	Rows       models.Rows
	URLBuilder *URLBuilder
}

// LoadTemplates loads HTML template files
func LoadTemplate(r *http.Request) (*template.Template, *TemplateData, error) {
	data := &TemplateData{
		URLBuilder: &URLBuilder{
			URL: r.URL,
		},
	}
	template, err := template.ParseFiles(baseTmpl, chartTmpl, dataTmpl)
	if err != nil {
		return nil, nil, err
	}
	return template, data, nil
}
