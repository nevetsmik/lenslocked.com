package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"lenslocked.com/controllers"
	"lenslocked.com/dbConfig"
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

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery)

	r := mux.NewRouter()
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
	r.Handle("/galleries/new", galleriesC.New).Methods("GET")
	r.HandleFunc("/galleries", galleriesC.Create).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")


	http.ListenAndServe(":3000", r)
}
