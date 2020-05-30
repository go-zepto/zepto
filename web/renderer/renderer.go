package renderer

import (
	"errors"
	"fmt"
	"github.com/flosch/pongo2"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// A Tmpl implements keeper, loader and reloader for HTML templates
type Tmpl struct {
	dir      string                      // root directory
	ext      string                      // extension
	watch    bool                        // reload every time
	loadedAt time.Time                   // loaded at (last loading time)
	tmplMap  map[string]*pongo2.Template // Map of key => Template
}

// NewTmpl creates new Tmpl and loads templates. The dir argument is
// directory to load templates from. The ext argument is extension of
// tempaltes. The watch (if true) turns the Tmpl to reload templates
// every Render if there is a change in the dir.
func NewTmpl(dir, ext string, watch bool) (tmpl *Tmpl, err error) {

	// get absolute path
	if dir, err = filepath.Abs(dir); err != nil {
		return
	}
	fmt.Println(dir)

	tmpl = new(Tmpl)
	tmpl.dir = dir
	tmpl.ext = ext
	tmpl.watch = watch
	tmpl.tmplMap = make(map[string]*pongo2.Template)

	if err = tmpl.Load(); err != nil {
		tmpl = nil // drop for GC
	}

	return
}

// Dir returns absolute path to directory with views
func (t *Tmpl) Dir() string {
	return t.dir
}

// Ext returns extension of views
func (t *Tmpl) Ext() string {
	return t.ext
}

// Devel returns watchopment pin
func (t *Tmpl) Devel() bool {
	return t.watch
}

func recoverTemplateNotFound() {
	if r := recover(); r != nil {
		err := r.(*pongo2.Error)
		if err.OrigError.Error() == "unable to resolve template" {
			fmt.Println("[pongo2] Unable to resolve template: " + err.Filename)
			os.Exit(1)
		}
		panic(r)
	}
}

// Load or reload templates
func (t *Tmpl) Load() (err error) {

	defer recoverTemplateNotFound()

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
		if filepath.Ext(path) != t.ext {
			return
		}

		// get relative path
		var rel string
		if rel, err = filepath.Rel(t.dir, path); err != nil {
			return err
		}

		// name of a template is its relative path
		// without extension
		rel = strings.TrimSuffix(rel, t.ext)
		tplExample := pongo2.Must(pongo2.FromFile(path))
		t.tmplMap[rel] = tplExample
		return err
	}

	if err = filepath.Walk(t.dir, walkFunc); err != nil {
		return
	}

	return
}

// IsModified lookups directory for changes to
// reload (or not to reload) templates if watchopment
// pin is true.
func (t *Tmpl) IsModified() (yep bool, err error) {

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
		if filepath.Ext(path) != t.ext {
			return
		}

		if yep = info.ModTime().After(t.loadedAt); yep == true {
			return errStop
		}

		return
	}

	// clear the errStop
	if err = filepath.Walk(t.dir, walkFunc); err == errStop {
		err = nil
	}

	return
}

func (t *Tmpl) Render(w io.Writer, name string, data pongo2.Context) (err error) {

	// if development (watch mode)
	if t.watch == true {

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
		fmt.Println("Error. Template not found: ", name)
	}

	tmpl.ExecuteWriter(data, w)
	return
}
