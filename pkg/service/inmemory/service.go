package inmemory

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/jlordiales/go-kit-todo-backend/pkg/service"
	"github.com/satori/go.uuid"
)

func NewService(logger log.Logger) service.Service {
	return &inMemoryService{make(map[string]service.Todo, 0), logger}
}

type inMemoryService struct {
	todos  map[string]service.Todo
	logger log.Logger
}

func (s *inMemoryService) Create(ctx context.Context, title string, order int) service.Todo {
	todo := service.TodoFrom(title, order)
	s.todos[todo.Id.String()] = todo
	return todo
}

func (s *inMemoryService) UpdateOrder(ctx context.Context, id uuid.UUID, order int) (service.Todo, error) {
	return updateTodo(ctx, s, id, func(todo service.Todo) service.Todo {
		todo.Order = order
		return todo
	})
}

func (s *inMemoryService) UpdateTitle(ctx context.Context, id uuid.UUID, title string) (service.Todo, error) {
	return updateTodo(ctx, s, id, func(todo service.Todo) service.Todo {
		todo.Title = title
		return todo
	})
}

func (s *inMemoryService) Complete(ctx context.Context, id uuid.UUID) (service.Todo, error) {
	return updateTodo(ctx, s, id, func(todo service.Todo) service.Todo {
		todo.Completed = true
		return todo
	})
}

func (s *inMemoryService) DeleteAllTodos(ctx context.Context) {
	s.todos = make(map[string]service.Todo, 0)
}

func (s *inMemoryService) DeleteSpecificTodo(ctx context.Context, id uuid.UUID) error {
	_, found := s.todos[id.String()]
	if found {
		delete(s.todos, id.String())
		return nil
	}
	return service.ErrTodoNotFound
}

func (s *inMemoryService) ListAllTodos(ctx context.Context) ([]service.Todo, error) {
	v := make([]service.Todo, 0, len(s.todos))

	for _, value := range s.todos {
		v = append(v, value)
	}
	return v, nil
}

func (s *inMemoryService) GetTodo(ctx context.Context, id uuid.UUID) (service.Todo, error) {
	var value service.Todo
	key := id.String()
	value, ok := s.todos[key]
	if !ok {
		s.logger.Log("id", id, "msg", "todo not found")
		return value, service.ErrTodoNotFound
	}
	return value, nil
}

func updateTodo(ctx context.Context, s *inMemoryService, id uuid.UUID, updateFunction func(service.Todo) service.Todo) (service.Todo, error) {
	var existingTodo service.Todo
	existingTodo, e := s.GetTodo(ctx, id)
	if e != nil {
		return existingTodo, e
	}

	updatedTodo := updateFunction(existingTodo)
	s.todos[id.String()] = updatedTodo
	return updatedTodo, nil
}
