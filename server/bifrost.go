// Heimdal! Open the Bifrost!
package server

import (
	"context"
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

	router, eventCh := handlers.NewRouter()
	eng.Process(context.Background(), eventCh)

	app, use := middleware.Wrap(router)
	use(middleware.NewCors("*"))
	use(&middleware.RequestLogger{})

	return http.ListenAndServe(":2727", app)
}
