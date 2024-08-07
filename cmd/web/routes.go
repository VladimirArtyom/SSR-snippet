package main

import ("net/http"
  "github.com/justinas/alice"
  "github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {

  var router *httprouter.Router = httprouter.New()
  
  router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    app.notFound(w) 
  })

  var fileServer http.Handler = http.FileServer(http.Dir("./ui/static/"))
  router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
  // Nous interessons uniqument a la session sur les routes ci-dessous.
  // We only interested the session on the routes below.
  sessionMiddleWare := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
  authMiddleware := sessionMiddleWare.Append(app.requireAuth)

  router.Handler(http.MethodGet, "/", sessionMiddleWare.ThenFunc(app.Home))
  router.Handler(http.MethodGet, "/snip/view/:id", sessionMiddleWare.ThenFunc(app.SnipView))
  
  // Authentications
  router.Handler(http.MethodGet, "/user/signup", sessionMiddleWare.ThenFunc(app.UserSignup))
  router.Handler(http.MethodPost, "/user/signup", sessionMiddleWare.ThenFunc(app.UserSignupPost))

  router.Handler(http.MethodGet, "/user/login", sessionMiddleWare.ThenFunc(app.UserLogin))
  router.Handler(http.MethodPost, "/user/login", sessionMiddleWare.ThenFunc(app.UserLoginPost))


  // Authorization & Authentications
  router.Handler(http.MethodGet, "/snip/create", authMiddleware.ThenFunc(app.SnipCreate))
  router.Handler(http.MethodPost, "/snip/create", authMiddleware.ThenFunc(app.SnipCreatePost))
  router.Handler(http.MethodPost, "/user/logout", authMiddleware.ThenFunc(app.UserLogoutPost))


  var standardMiddleWare alice.Chain = alice.New(app.recoverPanic, app.logRequest, secureHeaders)
  return standardMiddleWare.Then(router)
}
