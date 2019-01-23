// Heimdal! Open the Bifrost!
package bifrost

import (
	"net/http"

	"github.com/goslang/bifrost/handlers"
	"github.com/goslang/bifrost/lib/middleware"
)

func Start() error {
	router := handlers.NewRouter()

	app, use := middleware.Wrap(router)
	use(middleware.NewCors("*"))
	use(&middleware.RequestLogger{})

	return http.ListenAndServe(":2727", app)
}
