package renderer

import (
	"fmt"
	"html/template"
	"net/http"
	"runtime/debug"
	"strings"
)

type DevErrorData struct {
	Title string
	Res   http.ResponseWriter
	Req   *http.Request
	URI   string
	Trace template.HTML
}

func RenderDevelopmentError(w http.ResponseWriter, r *http.Request, err error) {
	t, _ := template.New("error").Parse(ERROR_HTML_TEMPLATE)
	trace := string(debug.Stack())
	trace = strings.Replace(trace, "\n", "<br>", -1)
	data := DevErrorData{
		Title: err.Error(),
		Res:   w,
		Req:   r,
		URI:   fmt.Sprintf("http://%s", r.Host),
		Trace: template.HTML(trace),
	}
	t.Execute(w, &data)
}
