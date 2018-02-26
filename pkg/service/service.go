package service

import (
	"context"
	"errors"

	"github.com/satori/go.uuid"
)

type Service interface {
	Create(ctx context.Context, title string, order int) Todo
	UpdateOrder(ctx context.Context, id uuid.UUID, order int) (Todo, error)
	UpdateTitle(ctx context.Context, id uuid.UUID, title string) (Todo, error)
	Complete(ctx context.Context, id uuid.UUID) (Todo, error)

	DeleteAllTodos(ctx context.Context)
	DeleteSpecificTodo(ctx context.Context, id uuid.UUID) error

	ListAllTodos(ctx context.Context) ([]Todo, error)
	GetTodo(ctx context.Context, id uuid.UUID) (Todo, error)
}

var ErrTodoNotFound = errors.New("value not found")
