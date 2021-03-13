package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"lenslocked.com/context"
	"lenslocked.com/email"
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
	NewView      *views.View
	LoginView    *views.View
	ForgotPwView *views.View
	ResetPwView  *views.View
	us           interfaces.UserServiceInt
	emailer      *email.Client
	r            *mux.Router
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

// ResetPwForm is used to process the forgot password form
// and the reset password form.
type ResetPwForm struct {
	Email    string `schema:"email"`
	Token    string `schema:"token"`
	Password string `schema:"password"`
}

type Alert struct {
	Level   string
	Message string
}

func NewUsers(us interfaces.UserServiceInt, r *mux.Router, emailer *email.Client) *Users {
	return &Users{
		NewView:      views.NewView("bootstrap", "users/new"),
		LoginView:    views.NewView("bootstrap", "users/login"),
		ForgotPwView: views.NewView("bootstrap", "users/forgot_pw"),
		ResetPwView:  views.NewView("bootstrap", "users/reset_pw"),
		us:           us,
		emailer:      emailer,
		r:            r,
	}
}

// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	// Parse url query params for form
	// E.g., http://localhost:3000/signup?email=michael%40dundermifflin.com
	parseURLParams(r, &form)
	u.NewView.Render(w, r, form)
}

// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form SignupForm
	vd.Yield = &form // Persist form data
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		// Errors are possible here. Display AlertMsgGeneric message
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
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
	alert := views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Welcome to LensLocked.com!",
	}
	// Redirect with alert
	views.RedirectAlert(w, r, url.Path, http.StatusFound, alert)
}

// POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)

	}
	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			vd.AlertError("No user exists with that email address")
		default:
			vd.SetAlert(err)
		}
		u.LoginView.Render(w, r, vd)
		return
	}
	err = u.signIn(w, user)
	if err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
	}
	url, err := u.r.Get(IndexGalleries).URL()
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}

// Logout is used to delete a user's session cookie and invalidate their current remember token, which will sign the current user out.
// POST /logout
func (u *Users) Logout(w http.ResponseWriter, r *http.Request) {
	// First expire the user's cookie
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	// Then we update the user with a new remember token
	user := context.User(r.Context())
	// We are ignoring errors for now because they are unlikely, and even if they do occur we can't recover
	// now that the user doesn't have a valid cookie
	token, _ := rand.RememberToken()
	user.Remember = token
	u.us.Update(user)
	// Finally send the user to the home page
	http.Redirect(w, r, "/", http.StatusFound)
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

// POST /forgot
func (u *Users) InitiateReset(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ResetPwForm
	vd.Yield = &form
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.ForgotPwView.Render(w, r, vd)
		return
	}

	token, err := u.us.InitiateReset(form.Email)
	if err != nil {
		vd.SetAlert(err)
		u.ForgotPwView.Render(w, r, vd)
		return
	}

	// Email the user their password reset token. In the
	err = u.emailer.ResetPw(form.Email, token)
	if err != nil {
		vd.SetAlert(err)
		u.ForgotPwView.Render(w, r, vd)
		return
	}

	views.RedirectAlert(w, r, "/reset", http.StatusFound, views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Instructions for resetting your password have been emailed to you.",
	})
}

// ResetPw displays the reset password form and has a method so that we can prefill the form data with a token provided
// via the URL query params.
//
// GET /reset
func (u *Users) ResetPw(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ResetPwForm
	vd.Yield = &form
	if err := parseURLParams(r, &form); err != nil {
		vd.SetAlert(err)
	}
	u.ResetPwView.Render(w, r, vd)
}

// CompleteReset processed the reset password form
//
// POST /reset
func (u *Users) CompleteReset(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ResetPwForm
	vd.Yield = &form
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.ResetPwView.Render(w, r, vd)
		return
	}

	user, err := u.us.CompleteReset(form.Token, form.Password)
	if err != nil {
		vd.SetAlert(err)
		u.ResetPwView.Render(w, r, vd)
		return
	}

	u.signIn(w, user)
	views.RedirectAlert(w, r, "/galleries", http.StatusFound, views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Your password has been reset and you have been logged in!",
	})
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
