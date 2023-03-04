package handler

import (
	"encoding/json"
	"net/http"
	"sync"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/shusann01116/Chatting/backend/app/errors"
	"github.com/shusann01116/Chatting/backend/app/loader"
)

// The GraphQL handler handles GraphQL API requests over HTTP.
// It can handle batched requests as sent by the apollo-client.
type GraphQL struct {
	Schema  *graphql.Schema
	Loaders loader.Collection
	Logger  logger
}

func (h GraphQL) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Validate the request.
	if ok := isSupported(r.Method); !ok {
		respond(w, errorJSON("unsupported method"), http.StatusMethodNotAllowed)
		return
	}

	req, err := parse(r)
	if err != nil {
		respond(w, errorJSON(err.Error()), http.StatusBadRequest)
		return
	}

	n := len(req.queries)
	if n == 0 {
		respond(w, errorJSON("no query found"), http.StatusBadRequest)
		return
	}

	// Authentication happens here

	// Execute the queries.
	var (
		ctx       = h.Loaders.Attach(r.Context())
		responses = make([]*graphql.Response, n)
		wg        sync.WaitGroup
	)

	wg.Add(n)

	for i, q := range req.queries {
		// Loop through the parased queries from the request.
		// These queries are executed in separate goroutines so they process in parallel.
		go func(i int, q query) {
			res := h.Schema.Exec(ctx, q.Query, q.OpName, q.Variables)

			// We have to do some work here to expand errors when it is possible for a resolver to return
			// more than one error (for example, a list resolver).
			res.Errors = errors.Expand(res.Errors)

			responses[i] = res
			wg.Done()
		}(i, q)
	}

	wg.Wait()

	var resp []byte
	if req.isBatch {
		resp, err = json.Marshal(responses)
	} else if len(responses) > 0 {
		resp, err = json.Marshal(responses[0])
	}

	if err != nil {
		respond(w, errorJSON("server error"), http.StatusInternalServerError)
	}

	respond(w, resp, http.StatusOK)
}

// logger defines an interface with a single method.
type logger interface {
	Printf(fmt string, values ...interface{})
}

// A request represents an HTTP request to the GraphQL endpoint.
// A request can have a single query or a batch of requests with one or more queries.
type request struct {
	queries []query
	isBatch bool
}

// A query represents a single GraphQL query.
type query struct {
	OpName    string                 `json:"operationName"`
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}
