package renderer

import (
	"fmt"
	"html/template"
	"net/http"
)

type DevErrorData struct {
	Title string
	Res   http.ResponseWriter
	Req   *http.Request
	URI   string
}

func RenderDevelopmentError(w http.ResponseWriter, r *http.Request, err error) {
	t, _ := template.New("error").Parse(ERROR_HTML_TEMPLATE)
	data := DevErrorData{
		Title: err.Error(),
		Res:   w,
		Req:   r,
		URI:   fmt.Sprintf("http://%s", r.Host),
	}
	t.Execute(w, &data)
}
