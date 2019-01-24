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
	eventsCh := make(chan engine.Event)

	router := handlers.NewRouter(eventsCh)

	app, use := middleware.Wrap(router)
	use(middleware.NewCors("*"))
	use(&middleware.RequestLogger{})

	go eng.Process(context.Background(), eventsCh)
	return http.ListenAndServe(":2727", app)
}
