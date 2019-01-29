// Heimdal! Open the Bifrost!
package bifrost

import (
	"context"
	"encoding/gob"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

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

	ctx := context.Background()
	go eng.Process(ctx, eventsCh)
	//go startSnapshots(ctx, eventsCh)
	return http.ListenAndServe(":2727", app)
}

func startSnapshots(ctx context.Context, eventsCh chan engine.Event) {
	timer := engine.SnapshotTimer(
		ctx,
		10*time.Second,
		encoderFactory,
		writerFactory,
	)

	for snapshotEvt := range timer {
		eventsCh <- snapshotEvt
	}
}

func encoderFactory(writer io.Writer) engine.Encoder {
	return gob.NewEncoder(writer)
}

func writerFactory() io.WriteCloser {
	filePrefix := "/usr/local/var/bifrost/"
	fileName := filePrefix + "data-" + strconv.Itoa(int(time.Now().Unix())) + ".bin"

	// TODO: Fix ignoring error
	file, _ := os.Create(fileName)
	return file
}
