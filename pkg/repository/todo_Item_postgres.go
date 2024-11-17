package repository

import (
	"fmt"
	"log"

	todo "github.com/fancurson/toDoList"
	"github.com/jmoiron/sqlx"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{
		db: db,
	}
}

func (r *TodoItemPostgres) CreateItem(listId int, item todo.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int

	createItemQuery := fmt.Sprintf(`
	INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id
	`, todoItemTable)

	err = tx.QueryRow(createItemQuery, item.Title, item.Description).Scan(&itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createListItemQuery := fmt.Sprintf(`
	INSERT INTO %s (list_id, item_id) VALUES ($1, $2) RETURNING id
	`, listItemTable)

	_, err = tx.Exec(createListItemQuery, listId, itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemId, tx.Commit()
}

func (r *TodoItemPostgres) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	items := make([]todo.TodoItem, 0)

	log.Println("check point 1")

	query := fmt.Sprintf(`
	SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti 
	INNER JOIN %s li on li.item_id = ti.id 
	INNER JOIN %s ul on ul.list_id = li.list_id
	WHERE li.list_id = $1 AND ul.user_id = $2 
	`, todoItemTable, listItemTable, userListTable)

	if err := r.db.Select(&items, query, listId, userId); err != nil {
		return nil, err
	}

	return items, nil
}
