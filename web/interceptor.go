package web

import (
	"errors"
	"github.com/go-zepto/zepto/web/renderer"
	"net/http"
	"strconv"
	"strings"
)

type Interceptor struct {
	env        string
	origWriter http.ResponseWriter
	origReq    *http.Request
	overridden bool
}

func (i *Interceptor) WriteHeader(rc int) {
	if rc != 500 && rc >= 400 && i.env == "development" {
		errMessage := strings.ToLower(strconv.Itoa(rc) + " - " + http.StatusText(rc))
		i.origWriter.Header().Set("content-type", "text/html")
		i.origWriter.WriteHeader(rc)
		renderer.RenderDevelopmentError(i.origWriter, i.origReq, errors.New(errMessage))
		i.overridden = true
	} else {
		i.origWriter.WriteHeader(rc)
	}
	// TODO: Add support for custom errors (production mode)
}

func (i *Interceptor) Write(b []byte) (int, error) {
	if !i.overridden {
		return i.origWriter.Write(b)
	}

	// Return nothing if we've overriden the response.
	return 0, nil
}

func (i *Interceptor) Header() http.Header {
	return i.origWriter.Header()
}

func ErrorHandler(h http.Handler, env string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w = &Interceptor{origWriter: w, origReq: r, env: env}
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
