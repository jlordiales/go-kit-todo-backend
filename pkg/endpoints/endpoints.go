package endpoints

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/jlordiales/go-kit-todo-backend/pkg/service"
	"github.com/satori/go.uuid"
)

type CreateTodoRequest struct {
	Title string
	Order int
}

type UpdateTodoRequest struct {
	Id        string
	Title     *string
	Order     *int
	Completed *bool
}

type TodoResponse struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Order     int    `json:"order"`
	Completed bool   `json:"completed"`
	Url       string `json:"url"`
}

type GetTodoRequest struct {
	Id string
}

func serviceToEndpointResponse(todo service.Todo, basePath string) TodoResponse {
	return TodoResponse{
		todo.Id.String(),
		todo.Title,
		todo.Order,
		todo.Completed,
		fmt.Sprintf("%s/%s", basePath, todo.Id.String())}
}

func makeUpdateTodoEndpoint(s service.Service, basePath string) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateTodoRequest)
		id, e := uuid.FromString(req.Id)
		if e != nil {
			return nil, e
		}

		if req.Completed != nil && *req.Completed {
			s.Complete(ctx, id)
		}

		if req.Title != nil {
			s.UpdateTitle(ctx, id, *req.Title)
		}

		if req.Order != nil {
			s.UpdateOrder(ctx, id, *req.Order)
		}

		todo, e := s.GetTodo(ctx, id)
		if e != nil {
			return nil, e
		}

		return serviceToEndpointResponse(todo, basePath), nil
	}
}

func makeCreateTodoEndpoint(s service.Service, basePath string) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateTodoRequest)
		createdTodo := s.Create(ctx, req.Title, req.Order)
		return serviceToEndpointResponse(createdTodo, basePath), nil
	}
}

func makeListAllTodosEndpoint(s service.Service, basePath string) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		todos, e := s.ListAllTodos(ctx)
		if e != nil {
			return nil, e
		}

		responseTodos := make([]TodoResponse, 0, len(todos))
		for _, v := range todos {
			responseTodos = append(responseTodos, serviceToEndpointResponse(v, basePath))
		}
		return responseTodos, nil
	}
}

func makeGetTodoEndpoint(s service.Service, basePath string) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetTodoRequest)
		id, e := uuid.FromString(req.Id)
		if e != nil {
			return nil, e
		}

		todo, e := s.GetTodo(ctx, id)
		if e != nil {
			return nil, e
		}

		return serviceToEndpointResponse(todo, basePath), nil
	}
}

func makeDeleteAllTodosEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		s.DeleteAllTodos(ctx)
		return nil, nil
	}
}

func makeDeleteTodoEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		id, e := uuid.FromString(request.(string))
		if e != nil {
			return nil, e
		}

		return nil, s.DeleteSpecificTodo(ctx, id)
	}
}

func New(s service.Service, basePath string, logger log.Logger) Endpoints {
	var createEndpoint endpoint.Endpoint
	{
		createEndpoint = makeCreateTodoEndpoint(s, basePath)
	}

	var listAllEndpoint endpoint.Endpoint
	{
		listAllEndpoint = makeListAllTodosEndpoint(s, basePath)
	}

	var getTodoEndpoint endpoint.Endpoint
	{
		getTodoEndpoint = makeGetTodoEndpoint(s, basePath)
	}

	var deleteAllTodosEndpoint endpoint.Endpoint
	{
		deleteAllTodosEndpoint = makeDeleteAllTodosEndpoint(s)
	}

	var deleteTodo endpoint.Endpoint
	{
		deleteTodo = makeDeleteTodoEndpoint(s)
	}

	var updateTodo endpoint.Endpoint
	{
		updateTodo = makeUpdateTodoEndpoint(s, basePath)
	}

	return Endpoints{
		CreateTodoEndpoint: createEndpoint,
		ListAllTodos:       listAllEndpoint,
		GetTodo:            getTodoEndpoint,
		DeleteAll:          deleteAllTodosEndpoint,
		DeleteTodo:         deleteTodo,
		UpdateTodo:         updateTodo,
	}
}

type Endpoints struct {
	CreateTodoEndpoint endpoint.Endpoint
	ListAllTodos       endpoint.Endpoint
	GetTodo            endpoint.Endpoint
	DeleteAll          endpoint.Endpoint
	DeleteTodo         endpoint.Endpoint
	UpdateTodo         endpoint.Endpoint
}
