package views

import "html/template"

type View struct {
	Template *template.Template
}

// files is a variadic parameter that can accept 0...n arguments; files is treated as a slice in the NewView function
// NewView(someSlice...) to unravel the items in slice when invoking a variadic function
func NewView(files ...string) *View {
	files = append(files, "views/layouts/footer.gohtml")
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
	}
}
