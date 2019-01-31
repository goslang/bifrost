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
	eventCh := make(chan engine.Event)

	router := handlers.NewRouter(eventCh)

	app, use := middleware.Wrap(router)
	use(middleware.NewCors("*"))
	use(&middleware.RequestLogger{})

	ctx := context.Background()
	go eng.Process(ctx, eventCh)
	return http.ListenAndServe(":2727", app)
}
