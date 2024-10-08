package ui

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/dskart/waterfall-engine/app"
	"github.com/dskart/waterfall-engine/ui/pkg/middleware"
	"github.com/dskart/waterfall-engine/ui/pkg/router"
)

var (
	//go:embed public
	publicFs     embed.FS
	publicServer http.FileSystem
)

func init() {
	fsys, err := fs.Sub(publicFs, "public")
	if err != nil {
		panic(err)
	}
	publicServer = http.FS(fsys)
}

const MaxProxyCount = 1

func RegisterHandles(r *mux.Router, app *app.App) error {
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(publicServer)))

	anonymousRouter := r.NewRoute().Subrouter()
	sm := middleware.NewSessionMiddleware(app, MaxProxyCount)
	anonymousRouter.Use(sm.AnonymousSession)
	router.InitComponentRouter(anonymousRouter)
	router.InitPageRouter(anonymousRouter)
	return nil
}
