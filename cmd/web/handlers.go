package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Jackob2004/snippetbox/internal/models"
	"github.com/Jackob2004/snippetbox/internal/validator"
)

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (a *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := a.snippets.Latest()

	if err != nil {
		a.serverError(w, r, err)
		return
	}

	userId := a.sessionManager.GetInt(r.Context(), AuthUserId)
	data := a.newTemplateData(r)
	data.Snippets = snippets
	data.User = models.User{ID: userId}
	a.render(w, r, http.StatusOK, "home.gohtml", data)
}

func (a *application) about(w http.ResponseWriter, r *http.Request) {
	data := a.newTemplateData(r)
	a.render(w, r, http.StatusOK, "about.gohtml", data)
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

	userId := a.sessionManager.GetInt(r.Context(), AuthUserId)
	data := a.newTemplateData(r)
	data.Snippet = snippet
	data.User = models.User{ID: userId}

	a.render(w, r, http.StatusOK, "view.gohtml", data)
}

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (f *snippetCreateForm) validate() bool {
	f.CheckField(validator.NotBlank(f.Title), "title", "This field cannot be blank")
	f.CheckField(validator.MaxChars(f.Title, 100), "title", "This field cannot be more than 100 characters long")
	f.CheckField(validator.NotBlank(f.Content), "content", "This field cannot be blank")
	f.CheckField(validator.PermittedValue(f.Expires, 1, 7, 365), "expires", "This field must be equal 1, 7 or 365")

	return f.Valid()
}

func (a *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := a.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}

	a.render(w, r, http.StatusOK, "create.gohtml", data)
}

func (a *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm

	err := a.decodePostForm(r, &form)
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
	}

	if !form.validate() {
		data := a.newTemplateData(r)
		data.Form = form
		a.render(w, r, http.StatusUnprocessableEntity, "create.gohtml", data)
		return
	}

	userId := a.sessionManager.GetInt(r.Context(), AuthUserId)
	id, err := a.snippets.Insert(form.Title, form.Content, form.Expires, userId)
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	a.sessionManager.Put(r.Context(), Flash, "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (a *application) snippetDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	err = a.snippets.Delete(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		}
		a.serverError(w, r, err)
		return
	}

	a.sessionManager.Put(r.Context(), Flash, "Snippet successfully deleted!")

	http.Redirect(w, r, "/account/snippets", http.StatusSeeOther)
}

func (a *application) snippetEdit(w http.ResponseWriter, r *http.Request) {
	snippetId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || snippetId < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := a.snippets.Get(snippetId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		}
		a.serverError(w, r, err)
		return
	}

	userId := a.sessionManager.GetInt(r.Context(), AuthUserId)
	if snippet.Creator.UserID != userId {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	data := a.newTemplateData(r)
	data.Snippet = models.Snippet{ID: snippetId}
	data.Form = snippetCreateForm{
		Title:   snippet.Title,
		Content: snippet.Content,
	}

	a.render(w, r, http.StatusOK, "edit.gohtml", data)
}

func (a *application) snippetEditPost(w http.ResponseWriter, r *http.Request) {
	snippetId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || snippetId < 1 {
		http.NotFound(w, r)
		return
	}

	var form snippetCreateForm

	err = a.decodePostForm(r, &form)
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
	}

	if !form.validate() {
		data := a.newTemplateData(r)
		data.Form = form
		a.render(w, r, http.StatusUnprocessableEntity, "edit.gohtml", data)
		return
	}

	userId := a.sessionManager.GetInt(r.Context(), AuthUserId)
	err = a.snippets.Update(snippetId, userId, form.Title, form.Content, form.Expires)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		}
		a.serverError(w, r, err)
		return
	}

	a.sessionManager.Put(r.Context(), Flash, "Snippet successfully edited!")
	a.logger.Info("TEST")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", snippetId), http.StatusSeeOther)
}

func (a *application) accountView(w http.ResponseWriter, r *http.Request) {
	id := a.sessionManager.GetInt(r.Context(), AuthUserId)
	user, err := a.users.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		} else {
			a.serverError(w, r, err)
		}
		return
	}

	data := a.newTemplateData(r)
	data.User = user
	a.render(w, r, http.StatusOK, "account.gohtml", data)
}

func (a *application) accountSnippets(w http.ResponseWriter, r *http.Request) {
	userId := a.sessionManager.GetInt(r.Context(), AuthUserId)
	snippets, err := a.snippets.GetSnippets(userId)
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	data := a.newTemplateData(r)
	data.Snippets = snippets
	a.render(w, r, http.StatusOK, "user_snippets.gohtml", data)
}

type passwordChangeForm struct {
	CurrentPassword     string `form:"current_password"`
	NewPassword         string `form:"new_password"`
	ConfirmPassword     string `form:"confirm_password"`
	validator.Validator `form:"-"`
}

func (a *application) accountPasswordUpdate(w http.ResponseWriter, r *http.Request) {
	data := a.newTemplateData(r)
	data.Form = passwordChangeForm{}

	a.render(w, r, http.StatusOK, "password.gohtml", data)
}

func (a *application) accountPasswordUpdatePost(w http.ResponseWriter, r *http.Request) {
	var form passwordChangeForm

	err := a.decodePostForm(r, &form)
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.CurrentPassword), "current_password", "This field cannot be blank")
	form.CheckField(validator.MaxBytes(form.CurrentPassword, 72), "current_password", "This field must not be more than 72 bytes long")
	form.CheckField(validator.NotBlank(form.NewPassword), "new_password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.NewPassword, 8), "new_password", "This field must be at least 8 characters long")
	form.CheckField(validator.MaxBytes(form.NewPassword, 72), "new_password", "This field must not be more than 72 bytes long")
	form.CheckField(validator.NotBlank(form.ConfirmPassword), "confirm_password", "This field cannot be blank")
	form.CheckField(form.NewPassword == form.ConfirmPassword, "confirm_password", "This field must be the same as new password")

	if !form.Valid() {
		data := a.newTemplateData(r)
		data.Form = form
		a.render(w, r, http.StatusUnprocessableEntity, "password.gohtml", data)
		return
	}

	userId := a.sessionManager.GetInt(r.Context(), AuthUserId)
	err = a.users.PasswordUpdate(userId, form.CurrentPassword, form.NewPassword)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Validator.AddFieldError("current_password", "Current password value is invalid")
			data := a.newTemplateData(r)
			data.Form = form
			a.render(w, r, http.StatusUnprocessableEntity, "password.gohtml", data)
		} else {
			a.serverError(w, r, err)
		}

		return
	}

	a.sessionManager.Put(r.Context(), Flash, "Passwords successfully updated!")
	http.Redirect(w, r, "/account/view", http.StatusSeeOther)
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (a *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := a.newTemplateData(r)
	data.Form = userSignupForm{}
	a.render(w, r, http.StatusOK, "signup.gohtml", data)
}

func (a *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	err := a.decodePostForm(r, &form)
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")
	form.CheckField(validator.MaxBytes(form.Password, 72), "password", "This field must not be more than 72 bytes long")

	if !form.Valid() {
		data := a.newTemplateData(r)
		data.Form = form
		a.render(w, r, http.StatusUnprocessableEntity, "signup.gohtml", data)
		return
	}

	err = a.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already taken")
			data := a.newTemplateData(r)
			data.Form = form
			a.render(w, r, http.StatusUnprocessableEntity, "signup.gohtml", data)
		} else {
			a.serverError(w, r, err)
		}
		return
	}

	a.sessionManager.Put(r.Context(), Flash, "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (a *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := a.newTemplateData(r)
	data.Form = userLoginForm{}
	a.render(w, r, http.StatusOK, "login.gohtml", data)
}

func (a *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := a.decodePostForm(r, &form)
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MaxBytes(form.Password, 72), "password", "This field must not be more than 72 bytes long")

	if !form.Valid() {
		data := a.newTemplateData(r)
		data.Form = form
		a.render(w, r, http.StatusUnprocessableEntity, "login.gohtml", data)
		return
	}

	id, err := a.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := a.newTemplateData(r)
			data.Form = form
			a.render(w, r, http.StatusUnprocessableEntity, "login.gohtml", data)
		} else {
			a.serverError(w, r, err)
		}

		return
	}

	err = a.sessionManager.RenewToken(r.Context())
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	a.sessionManager.Put(r.Context(), AuthUserId, id)

	path := a.sessionManager.PopString(r.Context(), RedirectPathAfterLogin)
	if path == "" {
		path = "/snippet/create"
	}
	http.Redirect(w, r, path, http.StatusSeeOther)
}

func (a *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := a.sessionManager.RenewToken(r.Context())
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	a.sessionManager.Remove(r.Context(), AuthUserId)
	a.sessionManager.Put(r.Context(), Flash, "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
