package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"lenslocked.com/controllers"
	"lenslocked.com/dbConfig"
	"lenslocked.com/middleware"
	"lenslocked.com/services"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Dbname)

	services, err := services.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.AutoMigrate()
	//services.DestructiveReset()

	r := mux.NewRouter()
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery, r)

	requireUserMw := middleware.RequireUser{
		UserServiceInt: services.User,
	}

	// staticC returns a struct of View structs.
	// Handle takes a path, and a http.Handler object.
	// Since View has a ServeHTTP method and implements the http.Handler interfaces, so the View type can be passed as
	// as http.Handler object
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	newGallery := requireUserMw.Apply(galleriesC.New)
	createGallery := requireUserMw.ApplyFn(galleriesC.Create)
	editGallery := requireUserMw.ApplyFn(galleriesC.Edit)
	updateGallery := requireUserMw.ApplyFn(galleriesC.Update)
	r.Handle("/galleries/new", newGallery).Methods("GET")
	r.HandleFunc("/galleries", createGallery).Methods("POST")
	// Name the route controllers.ShowGallery
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").Name(controllers.ShowGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/edit", editGallery).Methods("GET")
	r.HandleFunc("/galleries/{id:[0-9]+}/update", updateGallery).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")

	http.ListenAndServe(":3000", r)
}
