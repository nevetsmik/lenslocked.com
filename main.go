package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"lenslocked.com/controllers"
)

func main() {
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers()

	r := mux.NewRouter()
	// staticC returns a struct of View structs.
	// Handle takes a path, and a http.Handler object.
	// Since View has a ServeHTTP method and implements the http.Handler interface, so the View type can be passed as
	// as http.Handler object
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")

	http.ListenAndServe(":3000", r)
}
