package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type View struct {
	Layout   string
	Template *template.Template
}

var (
	LayoutDir   string = "views/layouts/"
	TemplateExt string = ".gohtml"
)

// files is a variadic parameter that can accept 0...n arguments; files is treated as a slice in the NewView function
// ewView(someSlice...) to unravel the items in slice when invoking a variadic function
func NewView(layout string, files ...string) *View {
	files = append(files,
		layoutFiles()...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Layout:   layout,
		Template: t,
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
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}
