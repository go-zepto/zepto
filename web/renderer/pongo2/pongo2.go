package pongo2

import (
	"errors"
	"github.com/flosch/pongo2"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// A Pongo2Engine implements keeper, loader and reloader for HTML templates
type Pongo2Engine struct {
	opts     Options
	loadedAt time.Time                   // loaded at (last loading time)
	tmplMap  map[string]*pongo2.Template // Map of key => Template
}

func NewPongo2Engine(opts ...Option) (tmpl *Pongo2Engine) {
	options := newOptions(opts...)
	absDir, err := filepath.Abs(options.templateDir)
	if err != nil {
		panic(err)
	}
	options.templateDir = absDir
	return &Pongo2Engine{
		opts:    options,
		tmplMap: make(map[string]*pongo2.Template),
	}
}

func (e *Pongo2Engine) Init() error {
	return e.Load()
}

func recoverTemplateNotFound() error {
	if r := recover(); r != nil {
		err := r.(*pongo2.Error)
		if err.OrigError.Error() == "unable to resolve template" {
			return err.OrigError
		}
		return err
	}
	return nil
}

// Load or reload templates
func (t *Pongo2Engine) Load() (err error) {

	err = recoverTemplateNotFound()
	if err != nil {
		return err
	}

	// time point
	t.loadedAt = time.Now()

	// unnamed root template
	//var root = template.New("")

	var walkFunc = func(path string, info os.FileInfo, err error) (_ error) {

		// handle walking error if any
		if err != nil {
			return err
		}

		// skip all except regular files
		// TODO (kostyarin): follow symlinks
		if !info.Mode().IsRegular() {
			return
		}

		// filter by extension
		if filepath.Ext(path) != t.opts.ext {
			return
		}

		// get relative path
		var rel string
		if rel, err = filepath.Rel(t.opts.templateDir, path); err != nil {
			return err
		}

		// name of a template is its relative path
		// without extension
		rel = strings.TrimSuffix(rel, t.opts.ext)
		tplExample := pongo2.Must(pongo2.FromFile(path))
		t.tmplMap[rel] = tplExample
		return err
	}

	if err = filepath.Walk(t.opts.templateDir, walkFunc); err != nil {
		return
	}

	return
}

// IsModified lookups directory for changes to
// reload (or not to reload) templates if autoReloadopment
// pin is true.
func (t *Pongo2Engine) IsModified() (yep bool, err error) {
	var errStop = errors.New("stop")
	var walkFunc = func(path string, info os.FileInfo, err error) (_ error) {

		// handle walking error if any
		if err != nil {
			return err
		}

		// skip all except regular files
		// TODO: follow symlinks
		if !info.Mode().IsRegular() {
			return
		}

		// filter by extension
		if filepath.Ext(path) != t.opts.ext {
			return
		}

		if yep = info.ModTime().After(t.loadedAt); yep == true {
			return errStop
		}

		return
	}
	// clear the errStop
	if err = filepath.Walk(t.opts.templateDir, walkFunc); err == errStop {
		err = nil
	}
	return
}

func (t *Pongo2Engine) Render(w http.ResponseWriter, status int, name string, d *sync.Map) (err error) {

	// Decode *sync.Map => pongo2Context
	data := pongo2.Context{}
	d.Range(func(key, value interface{}) bool {
		data[key.(string)] = value
		return true
	})

	// if development (autoReload mode)
	if t.opts.autoReload == true {
		// lookup directory for changes
		var modified bool
		if modified, err = t.IsModified(); err != nil {
			return
		}
		// reload
		if modified == true {
			if err = t.Load(); err != nil {
				return
			}
		}
	}
	tmpl, ok := t.tmplMap[name]
	if !ok {
		return errors.New("template not found - " + name)
	}
	w.WriteHeader(status)
	err = tmpl.ExecuteWriter(data, w)
	return
}
