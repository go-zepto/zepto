package web

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-zepto/zepto/web/renderer"
)

type Interceptor struct {
	app        *App
	origWriter http.ResponseWriter
	origReq    *http.Request
	overridden bool
}

func (i *Interceptor) WriteHeader(rc int) {
	errMessage := strings.ToLower(strconv.Itoa(rc) + " - " + http.StatusText(rc))
	if rc == 500 {
		i.origWriter.WriteHeader(rc)
		if i.app.opts.env == "production" {
			// TODO: Render custom error prod
			i.origWriter.Write([]byte(errMessage))
		}
		return
	}
	/*
		intercepts all error statuses except error 500, which is handled especially
		with our development debugger that identifies the error message.
	*/
	if rc >= 400 {
		if i.app.opts.env == "development" {
			i.origWriter.Header().Set("content-type", "text/html")
			i.origWriter.WriteHeader(rc)
			renderer.RenderDevelopmentError(i.origWriter, i.origReq, errors.New(errMessage))
			i.overridden = true
		} else {
			i.origWriter.WriteHeader(rc)
			i.origWriter.Write([]byte(errMessage))
			i.overridden = true
		}
	} else {
		i.origWriter.WriteHeader(rc)
	}
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

func ErrorHandler(app *App) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w = &Interceptor{origWriter: w, origReq: r, app: app}
		app.muxRouter.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
