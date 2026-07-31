package main

import (
	"fmt"
	"net/http"
)

func (a *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.RequestURI
	)

	a.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (a *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (a *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := a.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template page %s does not exist", page)
		a.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)
	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		a.serverError(w, r, err)
	}
}
