package controllers

import (
	"fmt"
	"net/http"

	"lenslocked.com/interfaces"
	"lenslocked.com/models"
	"lenslocked.com/views"
)

type Galleries struct {
	New *views.View
	gs  interfaces.GalleryServiceInt
}

type GalleryForm struct {
	Title string `schema:"title"`
}

func NewGalleries(gs interfaces.GalleryServiceInt) *Galleries {
	return &Galleries{
		New: views.NewView("bootstrap", "galleries/new"),
		gs:  gs,
	}
}

func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}
	gallery := models.Gallery{
		Title: form.Title,
	}
	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}
	fmt.Fprintln(w, gallery)

}
