package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	dynamicMiddleware := alice.New(app.session.Enable)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/article/create", dynamicMiddleware.ThenFunc(app.createArticleForm))
	mux.Post("/article/create", dynamicMiddleware.ThenFunc(app.createArticle))
	mux.Get("/article/:id", dynamicMiddleware.ThenFunc(app.showArticle))
	mux.Get("/category/:category", dynamicMiddleware.ThenFunc(app.showCategoryArticles))
	mux.Get("/contacts", dynamicMiddleware.ThenFunc(app.contacts))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
