package repository

import (
	"fmt"
	"log"

	todo "github.com/fancurson/toDoList"
	"github.com/jmoiron/sqlx"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListService(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{
		db: db,
	}
}

func (r *TodoListPostgres) Create(userId int, list todo.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	CreateListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListTable)
	err = r.db.QueryRow(CreateListQuery, list.Title, list.Description).Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createUserListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", userListTable)
	_, err = tx.Exec(createUserListQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *TodoListPostgres) GetAll(userId int) ([]todo.TodoList, error) {
	var lists []todo.TodoList

	query := fmt.Sprintf(`
    	SELECT tl.id, tl.title, tl.description 
    	FROM %s tl 
    	INNER JOIN %s ul ON tl.id = ul.list_id 
    	WHERE ul.user_id = $1`, todoListTable, userListTable,
	)
	err := r.db.Select(&lists, query, userId)

	log.Printf("%+v", lists)

	return lists, err
}

func (r *TodoListPostgres) GetById(userId, id int) (todo.TodoList, error) {

	var list todo.TodoList

	query := fmt.Sprintf(`
    SELECT tl.id, tl.title, tl.description 
    FROM %s tl 
    INNER JOIN %s ul ON tl.id = ul.list_id 
    WHERE ul.user_id = $1 AND tl.id = $2`,
		todoListTable, userListTable,
	)
	err := r.db.Get(&list, query, userId, id)

	return list, err
}
