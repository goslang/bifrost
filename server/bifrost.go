// Heimdal! Open the Bifrost!
package server

import (
	"net/http"

	"github.com/goslang/bifrost/engine"
	"github.com/goslang/bifrost/server/handlers"
	"github.com/goslang/bifrost/server/lib/middleware"
)

func Start() error {
	eng := engine.New()

	go func() {
		if err := eng.Run(); err != nil {
			panic(err)
		}
	}()

	router := handlers.NewRouter(eng)

	app, use := middleware.Wrap(router)
	use(middleware.NewCors("*"))
	use(&middleware.RequestLogger{})

	return http.ListenAndServe(":2727", app)
}
