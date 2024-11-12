package handler

import (
	"fmt"
	"go-playground/cmd/api/internal/datasource"
	"go-playground/cmd/api/internal/datasource/database"
	"go-playground/cmd/api/internal/transportlayer/rest"
	"go-playground/cmd/api/internal/transportlayer/rest/middleware"
	"go-playground/cmd/api/internal/transportlayer/rest/oapi"
	"go-playground/cmd/api/internal/usecase"
	"go-playground/pkg/env/v2"
	"net/http"
	"time"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/tecchu11/nrgo-std/nrhttp"
)

type handlers struct {
	*HealthHandler
	*TaskHandler
	*UserHandler
}

// New creates handler to handle requests.
func New(app *newrelic.Application, lookup func(string) (string, bool)) (http.Handler, error) {
	db, queries, err := database.NewQueryDB(lookup)
	if err != nil {
		return nil, fmt.Errorf("new query db: %w", err)
	}
	transactionAdaptor := datasource.NewDBTransactionAdaptor(db)
	taskAdaptor := datasource.NewTaskAdaptor(queries)
	userAdaptor := datasource.NewUserAdaptor(queries)

	taskUseCase := usecase.NewTaskUseCase(taskAdaptor, transactionAdaptor)
	userUseCase := usecase.NewUserUseCase(userAdaptor)

	health := &HealthHandler{Pinger: db}
	task := &TaskHandler{TaskInteractor: taskUseCase}
	user := &UserHandler{UserInteractor: userUseCase}

	applier := env.New(lookup)
	issuer := applier.URL("AUTH_ISSUER_URL")
	if err := applier.Err(); err != nil {
		return nil, fmt.Errorf("find auth issuer url from env: %w", err)
	}
	roundTripper := http.DefaultTransport
	if t, ok := http.DefaultTransport.(*http.Transport); ok {
		clone := t.Clone()
		clone.IdleConnTimeout = 120 * time.Second
		clone.MaxIdleConns = 0
		clone.MaxIdleConnsPerHost = 10
		clone.MaxConnsPerHost = 100
		clone.TLSHandshakeTimeout = 5 * time.Second
		roundTripper = clone
	}
	auth, err := middleware.NewAuth(
		jwks.NewCachingProvider(issuer, 5*time.Minute, jwks.WithCustomClient(&http.Client{
			Transport: newrelic.NewRoundTripper(roundTripper),
			Timeout:   5 * time.Second,
		})),
		[]string{"account"},
		middleware.WithSkipRoute("GET /health"),
	)
	if err != nil {
		return nil, fmt.Errorf("new auth middleware: %w", err)
	}

	return oapi.HandlerWithOptions(
		&handlers{
			TaskHandler:   task,
			HealthHandler: health,
			UserHandler:   user,
		},
		oapi.StdHTTPServerOptions{
			Middlewares: []oapi.MiddlewareFunc{
				auth.CheckAccessToken,
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
