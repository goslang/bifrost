// Heimdal! Open the Bifrost!
package bifrost

import (
	"context"
	"net/http"

	"github.com/goslang/bifrost/engine"
	"github.com/goslang/bifrost/handlers"
	"github.com/goslang/bifrost/lib/middleware"
)

func Start() error {
	eng := engine.New()
	go eng.Run()

	router, eventCh := handlers.NewRouter()
	eng.Process(context.Background(), eventCh)

	app, use := middleware.Wrap(router)
	use(middleware.NewCors("*"))
	use(&middleware.RequestLogger{})

	return http.ListenAndServe(":2727", app)
}
