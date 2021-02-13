package views

import "html/template"

type View struct {
	Layout   string
	Template *template.Template
}

// files is a variadic parameter that can accept 0...n arguments; files is treated as a slice in the NewView function
// ewView(someSlice...) to unravel the items in slice when invoking a variadic function
func NewView(layout string, files ...string) *View {
	files = append(files,
		"views/layouts/bootstrap.gohtml",
		"views/layouts/footer.gohtml",
		"views/layouts/navbar.gohtml")
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Layout:   layout,
		Template: t,
	}
}
