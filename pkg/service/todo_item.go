package service

import (
	todo "github.com/fancurson/toDoList"
	"github.com/fancurson/toDoList/pkg/repository"
)

type TodoItemService struct {
	repo     repository.TodoItem
	repoList repository.TodoList
}

func NewTodoItemService(repo repository.TodoItem, repoList repository.TodoList) *TodoItemService {
	return &TodoItemService{
		repo:     repo,
		repoList: repoList,
	}
}

func (s TodoItemService) CreateItem(userId, listId int, item todo.TodoItem) (int, error) {
	_, err := s.repoList.GetById(userId, listId)
	if err != nil {
		return 0, err
	}
	return s.repo.CreateItem(listId, item)
}

func (s TodoItemService) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	_, err := s.repoList.GetById(userId, listId)
	if err != nil {
		return nil, err
	}
	return s.repo.GetAll(userId, listId)
}

func (s *TodoItemService) GetById(userId, itemId int) (todo.TodoItem, error) {
	return s.repo.GetById(userId, itemId)
}

func (s *TodoItemService) Delete(userId, itemId int) error {
	return s.repo.Delete(userId, itemId)
}

func (s *TodoItemService) Update(userId, itemId int, input todo.UpdateItemInput) error {
	return s.repo.Update(userId, itemId, input)
}
