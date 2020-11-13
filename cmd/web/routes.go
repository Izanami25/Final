package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, secureHeaders)
	// Use the nosurf middleware on all our 'dynamic' routes.
	dynamicMiddleware := alice.New(app.session.Enable, noSurf)
	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/post/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createPostForm))
	mux.Post("/post/create", dynamicMiddleware.Append(app.requireAuthentication, app.logRequest).ThenFunc(app.createPost))
	mux.Get("/post/:id", dynamicMiddleware.ThenFunc(app.showPost))
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))
	mux.Get("/user/contact", dynamicMiddleware.ThenFunc(app.contact))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return standardMiddleware.Then(mux)
}
