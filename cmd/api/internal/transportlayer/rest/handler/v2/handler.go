package handler

import (
	"fmt"
	"go-playground/cmd/api/internal/datasource"
	"go-playground/cmd/api/internal/datasource/maindb"
	"go-playground/cmd/api/internal/transportlayer/rest"
	"go-playground/cmd/api/internal/transportlayer/rest/middleware"
	"go-playground/cmd/api/internal/transportlayer/rest/oapi"
	"go-playground/cmd/api/internal/usecase"
	"net/http"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/tecchu11/nrgo-std/nrhttp"
)

type handlers struct {
	*HealthHandler
	*TaskHandler
}

// New creates handler to handle requests.
func New(app *newrelic.Application, lookup func(string) (string, bool)) (http.Handler, error) {
	db, queries, err := maindb.NewQueryDB(lookup)
	if err != nil {
		return nil, fmt.Errorf("new query db: %w", err)
	}
	transactionAdaptor := datasource.NewDBTransactionAdaptor(db)
	taskAdaptor := datasource.NewTaskAdaptor(queries)

	taskUseCase := usecase.NewTaskUseCase(taskAdaptor, transactionAdaptor)

	health := &HealthHandler{Pinger: db}
	task := &TaskHandler{TaskInteractor: taskUseCase}

	return oapi.HandlerWithOptions(
		&handlers{
			TaskHandler:   task,
			HealthHandler: health,
		},
		oapi.StdHTTPServerOptions{
			Middlewares: []oapi.MiddlewareFunc{
				middleware.Recover,
				nrhttp.Middleware(app),
			},
			ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
				txn := newrelic.FromContext(r.Context())
				txn.NoticeExpectedError(err)
				rest.Err(w, err.Error(), http.StatusBadRequest)
			},
		},
	), nil
}

var _ oapi.ServerInterface = (*handlers)(nil)
