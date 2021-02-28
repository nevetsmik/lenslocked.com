package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"lenslocked.com/controllers"
	"lenslocked.com/dbConfig"
	"lenslocked.com/middleware"
	"lenslocked.com/rand"
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
	usersC := controllers.NewUsers(services.User, r)
	galleriesC := controllers.NewGalleries(services.Gallery, services.Image, r)

	// Redirects to /login if a user is not signed in
	requireUserMw := middleware.RequireUser{}

	// Writes user to context if remember token is found
	// Moves to next(w, r) regardless
	userMw := middleware.User{UserServiceInt: services.User}

	isProd := false
	b, err := rand.Bytes(32)
	if err != nil {
		panic(err)
	}
	// csrfMw will check for a valid CSRF token any time a form is submitted or our server gets an HTTP POST web request
	csrfMw := csrf.Protect(b, csrf.Secure(isProd))

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
	deleteGallery := requireUserMw.ApplyFn(galleriesC.Delete)
	indexGallery := requireUserMw.ApplyFn(galleriesC.Index)
	uploadImage := requireUserMw.ApplyFn(galleriesC.ImageUpload)
	deleteImage := requireUserMw.ApplyFn(galleriesC.ImageDelete)
	r.Handle("/galleries/new", newGallery).Methods("GET")
	r.HandleFunc("/galleries", createGallery).Methods("POST")
	// Name the route controllers.ShowGallery
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").Name(controllers.ShowGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/edit", editGallery).Methods("GET").Name(controllers.EditGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/update", updateGallery).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", deleteGallery).Methods("POST")
	r.Handle("/galleries", indexGallery).Methods("GET").Name(controllers.IndexGalleries)
	r.HandleFunc("/galleries/{id:[0-9]+}/images", uploadImage).Methods("POST")
	imageHandler := http.FileServer(http.Dir("./images/"))
	// http.StripPrefix acts as middleware and removes "/images/" before passing to imageHandler
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", imageHandler))
	r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete", deleteImage).Methods("POST")

	assetHandler := http.FileServer(http.Dir("./assets/"))
	assetHandler = http.StripPrefix("/assets/", assetHandler)
	r.PathPrefix("/assets/").Handler(assetHandler)

	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	http.ListenAndServe(":3000", csrfMw(userMw.Apply(r)))
}
