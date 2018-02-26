package inmemory

import (
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/jlordiales/go-kit-todo-backend/pkg/service"
	"github.com/stretchr/testify/assert"
)

var s service.Service = NewService(log.NewNopLogger())

func TestAddTodo(t *testing.T) {
	todo := s.Create(nil, "My awesome title", 15)

	assert.Equal(t, "My awesome title", todo.Title)
	assert.Equal(t, 15, todo.Order)
	assert.NotNil(t, todo.Id)
}

func TestTodosAreCreatedAsNotCompleted(t *testing.T) {
	todo := givenAnExistingTodo()

	assert.False(t, todo.Completed)
}

func TestAfterAddingATodoItIsPossibleToGetItById(t *testing.T) {
	t1 := givenAnExistingTodo()

	t2, e := s.GetTodo(nil, t1.Id)

	assert.NoError(t, e)
	assert.Equal(t, t1, t2)
}

func TestItIsPossibleToListAllTodos(t *testing.T) {
	givenAnExistingTodo()
	givenAnExistingTodo()

	todos, e := s.ListAllTodos(nil)

	assert.NoError(t, e)
	assert.Len(t, todos, 2)
}

func TestTodosCanHaveTheirOrderUpdated(t *testing.T) {
	t1 := givenAnExistingTodo()

	updated, _ := s.UpdateOrder(nil, t1.Id, 99)
	t2, _ := s.GetTodo(nil, t1.Id)

	assert.Equal(t, updated, t2)
}

func TestTodosCanHaveTheirTitleUpdated(t *testing.T) {
	t1 := givenAnExistingTodo()

	updated, _ := s.UpdateTitle(nil, t1.Id, "New title")
	t2, _ := s.GetTodo(nil, t1.Id)

	assert.Equal(t, updated, t2)
}

func givenAnExistingTodo() service.Todo {
	return s.Create(nil, "My awesome title", 15)
}
