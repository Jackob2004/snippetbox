package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Jackob2004/snippetbox/internal/models"
)

func (a *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	snippets, err := a.snippets.Latest()

	if err != nil {
		a.serverError(w, r, err)
		return
	}

	data := a.newTemplateData(r)
	data.Snippets = snippets
	a.render(w, r, http.StatusOK, "home.gohtml", data)
}

func (a *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := a.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			a.serverError(w, r, err)
		}
		return
	}

	data := a.newTemplateData(r)
	data.Snippet = snippet

	a.render(w, r, http.StatusOK, "view.gohtml", data)
}

func (a *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

func (a *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := a.snippets.Insert(title, content, expires)
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
