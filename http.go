package zepto

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fatih/color"
)

type HealthHandler struct {
	z    *Zepto
	next http.Handler
}

type HealthStatus struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Age     string `json:"age"`
}

type HTTPZeptoHandler struct {
	z       *Zepto
	handler http.Handler
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (h *HTTPZeptoHandler) coloredMethod(method string) string {
	mcolor := color.New()
	switch method {
	case "GET":
		mcolor = color.New(color.FgCyan, color.Bold)
	case "POST":
		mcolor = color.New(color.FgGreen, color.Bold)
	case "PUT", "PATCH":
		mcolor = color.New(color.FgYellow, color.Bold)
	case "DELETE":
		mcolor = color.New(color.FgRed, color.Bold)
	case "OPTIONS", "HEAD":
		mcolor = color.New(color.FgBlue, color.Bold)
	}
	return mcolor.Sprint(method)
}

func (h *HTTPZeptoHandler) coloredStatus(status int) string {
	scolor := color.New()
	if status >= 200 && status <= 299 {
		scolor = color.New(color.FgGreen, color.Bold)
	} else if status >= 300 && status <= 499 {
		scolor = color.New(color.FgYellow, color.Bold)
	} else if status >= 500 {
		scolor = color.New(color.FgRed, color.Bold)
	}
	return scolor.Sprint(status)
}

func (h *HTTPZeptoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/health" {
		var age time.Duration
		if h.z.startedAt != nil {
			age = time.Since(*h.z.startedAt)
		}
		json.NewEncoder(w).Encode(HealthStatus{
			Name:    h.z.opts.Name,
			Version: h.z.opts.Version,
			Age:     age.Round(time.Second).String(),
		})
		return
	}
	t := time.Now()

	lrw := NewLoggingResponseWriter(w)
	h.handler.ServeHTTP(lrw, r)
	boldColor := color.New(color.Bold)
	method := h.coloredMethod(r.Method)
	took := boldColor.Sprint(time.Since(t).Round(time.Nanosecond).String())
	status := h.coloredStatus(lrw.statusCode)
	h.z.Logger().Infof("%s %s took=%s status=%s", method, r.URL.Path, took, status)
}
