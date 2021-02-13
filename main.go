package main

import (
	"lenslocked.com/views"
	"net/http"

	"github.com/gorilla/mux"
)

var homeView *views.View
var contactView *views.View

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	data := struct{ Name string }{Name: "Steve"}
	// Execute a template name homeView.Layout ("bootstrap") writing the results to w and passing data
	if err := homeView.Template.ExecuteTemplate(w, homeView.Layout, data); err != nil {
		panic(err)
	}
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := contactView.Template.ExecuteTemplate(w, contactView.Layout, nil); err != nil {
		panic(err)
	}
}

func main() {
	// Pass the name of the layout to use and specify the view to parse
	homeView = views.NewView("bootstrap","views/home.gohtml")
	contactView = views.NewView("bootstrap", "views/contact.gohtml")

	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)

	http.ListenAndServe(":3000", r)
}
