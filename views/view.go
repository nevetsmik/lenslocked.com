package views

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/gorilla/csrf"

	"lenslocked.com/context"
)

type View struct {
	Layout   string
	Template *template.Template
}

var (
	LayoutDir   string = "views/layouts/"
	TemplateDir string = "views/"
	TemplateExt string = ".gohtml"
)

// files is a variadic parameter that can accept 0...n arguments; files is treated as a slice in the NewView function
// NewView(someSlice...) to unravel the items in slice when invoking a variadic function
func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, layoutFiles()...)
	t, err := template.New("").Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			// If this is called without being replace with a proper implementation
			// returning an error as the second argument will cause our template
			// package to return an error when executed.
			return "", errors.New("csrfField is not implemented")
		},
		"pathEscape": func(s string) string {
			return url.PathEscape(s)
		},
		"isLoggedIn": func() bool {
			return false
		},
		// Once we have our template with a function we are going to pass in files
		// to parse, much like we were previously.
	}).ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

// Put all the "views/layouts/*.gohtml" files in a slice
func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

// Abstract away Template.ExecuteTemplate
func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	var vd Data
	switch d := data.(type) {
	case Data:
		vd = d
	default:
		vd = Data{
			Yield: data,
		}
	}

	// Lookup the alert and assign it if one is persisted
	if alert := getAlert(r); alert != nil {
		vd.Alert = alert
		clearAlert(w)
	}

	vd.User = context.User(r.Context())
	// Handle errors when rendering views by writing to a temp buffer and then copy to the response after the render
	// completes
	var buf bytes.Buffer
	// We need to create the csrfField using the current http request.
	csrfField := csrf.TemplateField(r)
	tpl := v.Template.Funcs(template.FuncMap{
		// We can also change the return type of our function, since we no longer
		// need to worry about errors.
		"csrfField": func() template.HTML {
			// We can then create this closure that returns the csrfField for
			// any templates that need access to it.
			return csrfField
		},
		"isLoggedIn": func() bool {
			if vd.User != nil {
				return true
			}
			return false
		},
	})
	// Then we continue to execute the template just like before.
	err := tpl.ExecuteTemplate(&buf, v.Layout, vd)
	if err != nil {
		http.Error(w, "Something went wrong. If the problem persists, please email support@lenslocked.com",
			http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil)
}

// addTemplatePath takes in a slice of strings
// representing file paths for templates, and it prepends
// the TemplateDir directory to each string in the slice
//
// Eg the input {"home"} would result in the output
// {"views/home"} if TemplateDir == "views/"
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

// addTemplateExt takes in a slice of strings
// representing file paths for templates and it appends
// the TemplateExt extension to each string in the slice
//
// Eg the input {"home"} would result in the output
// {"home.gohtml"} if TemplateExt == ".gohtml"
func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}
