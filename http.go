package zepto

import (
	"encoding/json"
	"net/http"
	"time"
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

func (h *HTTPZeptoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/health" {
		age := time.Since(*h.z.startedAt)
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
	took := time.Since(t).Round(time.Nanosecond).String()
	status := lrw.statusCode
	h.z.Logger().Infof("%s %s took=%s status=%d", r.Method, r.URL.Path, took, status)
}
