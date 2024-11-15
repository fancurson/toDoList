package repository

import (
	"fmt"
	"log"
	"strings"

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

	return lists, err
}

func (r *TodoListPostgres) GetById(userId, id int) (todo.TodoList, error) {

	var list todo.TodoList

	log.Printf("u: %d, i: %d", userId, id)
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

func (r *TodoListPostgres) Delete(userId, id int) error {

	query := fmt.Sprintf(`
    DELETE FROM %s tl 
	USING %s ul 
	WHERE tl.id = ul.list_id AND ul.user_id=$1 AND ul.list_id=$2`,
		todoListTable, userListTable,
	)
	_, err := r.db.Exec(query, userId, id)

	return err
}

func (r *TodoListPostgres) Update(userId, listId int, input todo.UpdateListInput) error {

	log.Printf("userID: %d", userId)
	log.Printf("listId: %d", listId)
	setValue := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		log.Println("Check1")
		setValue = append(setValue, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		log.Println("Check2")
		setValue = append(setValue, fmt.Sprintf("Description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	setQuery := strings.Join(setValue, ", ")

	query := fmt.Sprintf(`
	UPDATE %s tl
	SET %s
	FROM %s ul
	WHERE tl.id = ul.list_id AND ul.user_id=$%d AND ul.list_id=$%d
	`, todoListTable, setQuery, userListTable, argId, argId+1)

	args = append(args, userId, listId)

	log.Printf("updateQuery: %s", query)
	log.Printf("args: %s", args...)

	_, err := r.db.Exec(query, args...)
	return err
}
