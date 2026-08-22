package main

import (
	"net/http"

	"github.com/Jackob2004/snippetbox/ui"
	"github.com/justinas/alice"
)

func (a *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	mux.HandleFunc("GET /ping", ping)

	dynamic := alice.New(a.sessionManager.LoadAndSave, preventCSRF, a.authenticate)

	mux.Handle("GET /{$}", dynamic.ThenFunc(a.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(a.snippetView))
	mux.Handle("GET /about", dynamic.ThenFunc(a.about))

	mux.Handle("GET /user/signup", dynamic.ThenFunc(a.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(a.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(a.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(a.userLoginPost))

	protected := dynamic.Append(a.requireAuthentication)

	mux.Handle("POST /user/logout", protected.ThenFunc(a.userLogoutPost))
	mux.Handle("GET /snippet/create", protected.ThenFunc(a.snippetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(a.snippetCreatePost))
	mux.Handle("POST /snippet/delete/{id}", protected.ThenFunc(a.snippetDelete))
	mux.Handle("GET /snippet/edit/{id}", protected.ThenFunc(a.snippetEdit))
	mux.Handle("POST /snippet/edit/{id}", protected.ThenFunc(a.snippetEditPost))
	mux.Handle("GET /account/view", protected.ThenFunc(a.accountView))
	mux.Handle("GET /account/snippets", protected.ThenFunc(a.accountSnippets))
	mux.Handle("GET /account/password/update", protected.ThenFunc(a.accountPasswordUpdate))
	mux.Handle("POST /account/password/update", protected.ThenFunc(a.accountPasswordUpdatePost))

	standard := alice.New(a.recoverPanic, a.logRequest, commonHeaders)

	return standard.Then(mux)
}
