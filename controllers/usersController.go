package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"

	"lenslocked.com/interfaces"
	"lenslocked.com/models"
	"lenslocked.com/rand"
	"lenslocked.com/views"
)

var (
	LayoutDir   string = "views/layouts/"
	TemplateDir string = "views/"
	TemplateExt string = ".gohtml"
)

// Users is the type for the usersController
// The users model is accessed through the userService
type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        interfaces.UserServiceInt
	r         *mux.Router
}

// Struct tags are metadata that can be added to fields of any struct
// gorilla/schema package unmarshals JSON into a struct
type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type Alert struct {
	Level   string
	Message string
}

func NewUsers(us interfaces.UserServiceInt, r *mux.Router) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
		r:         r,
	}
}

// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, nil)
}

// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		// Errors are possible here. Display AlertMsgGeneric message
		vd.SetAlert(err)
		u.NewView.Render(w, vd)
		return
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, vd)
		return
	}
	err := u.signIn(w, &user)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	url, err := u.r.Get(IndexGalleries).URL()
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}

// POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, vd)

	}
	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			vd.AlertError("No user exists with that email address")
		default:
			vd.SetAlert(err)
		}
		u.LoginView.Render(w, vd)
		return
	}
	err = u.signIn(w, user)
	if err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, vd)
	}
	url, err := u.r.Get(IndexGalleries).URL()
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}

// signIn is used to sign the given user in via cookies
func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}

	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	return nil
}

func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := u.us.ByRemember(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, user)
}

func parseForm(r *http.Request, dst interface{}) error {
	// ParseForm parses a form on a request and stores the data in PostForm in the request
	if err := r.ParseForm(); err != nil {
		return err
	}
	dec := schema.NewDecoder()
	if err := dec.Decode(dst, r.PostForm); err != nil {
		return err
	}
	return nil
}
