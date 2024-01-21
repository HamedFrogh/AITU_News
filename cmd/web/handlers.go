package main

import (
	"errors"
	"fmt"
	"hamedfrogh.net/aitunews/pkg/forms"
	"hamedfrogh.net/aitunews/pkg/models"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s, err := app.articles.Latest(ctx)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Retrieve the list of categories
	categories, err := app.articles.GetCategories()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Articles:   s,
		Categories: categories,
	})
}

func (app *application) showArticle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.articles.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{
		Article: s,
	})
}

func (app *application) createArticleForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})

}

func (app *application) createArticle(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires", "category")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	id, err := app.articles.Insert(form.Get("title"), form.Get("content"), form.Get("expires"), form.Get("category"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Article successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/article/%d", id), http.StatusSeeOther)
}

func (app *application) showCategoryArticles(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get(":category")

	articles, err := app.articles.GetByCategory(category)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "category.page.tmpl", &templateData{
		Category: category,
		Articles: articles,
	})
}

func (app *application) contacts(w http.ResponseWriter, r *http.Request) {

	app.render(w, r, "contacts.page.tmpl", &templateData{})
}
