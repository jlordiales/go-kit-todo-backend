package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/jlordiales/go-kit-todo-backend/pkg/endpoints"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

func MakeHandler(logger log.Logger, endpoints endpoints.Endpoints) http.Handler {
	r := mux.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
	}

	r.Methods("POST").Path("/").Handler(httptransport.NewServer(
		endpoints.CreateTodoEndpoint,
		decodeCreateTodoRequest,
		httptransport.EncodeJSONResponse,
		options...,
	))

	r.Methods("GET").Path("/").Handler(httptransport.NewServer(
		endpoints.ListAllTodos,
		httptransport.NopRequestDecoder,
		httptransport.EncodeJSONResponse,
		options...,
	))

	r.Methods("GET").Path("/{id}").Handler(httptransport.NewServer(
		endpoints.GetTodo,
		decodeGetTodoRequest,
		httptransport.EncodeJSONResponse,
		options...,
	))

	r.Methods("DELETE").Path("/").Handler(httptransport.NewServer(
		endpoints.DeleteAll,
		httptransport.NopRequestDecoder,
		httptransport.EncodeJSONResponse,
		options...,
	))

	r.Methods("DELETE").Path("/{id}").Handler(httptransport.NewServer(
		endpoints.DeleteTodo,
		decodeDeleteTodoRequest,
		httptransport.EncodeJSONResponse,
		options...,
	))

	r.Methods("PATCH").Path("/{id}").Handler(httptransport.NewServer(
		endpoints.UpdateTodo,
		decodeUpdateTodoRequest,
		httptransport.EncodeJSONResponse,
		options...,
	))

	return r
}

func decodeCreateTodoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request endpoints.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeUpdateTodoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}

	var request endpoints.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	request.Id = id
	return request, nil
}

func decodeGetTodoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return endpoints.GetTodoRequest{Id: id}, nil
}

func decodeDeleteTodoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return id, nil
}
