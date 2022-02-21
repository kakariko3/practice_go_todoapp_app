package models

import (
	"log"
	"time"
)

type Todo struct {
	ID        int
	Content   string
	UserID    int
	CreatedAt time.Time
}

// Todoを作成 (Create)
func (u *User) CreateTodo(content string) (err error) {
	cmd := `INSERT INTO todos (
		content,
		user_id,
		created_at
	) VALUES ($1, $2, $3)`
	_, err = Db.Exec(cmd, content, u.ID, time.Now())
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

// Todoの取得 (Read) (1件取得)
func GetTodo(id int) (todo Todo, err error) {
	cmd := `SELECT id, content, user_id, created_at FROM todos WHERE id = $1`
	todo = Todo{} // 名前付き戻り値で初期化済み -> 不要かも?
	err = Db.QueryRow(cmd, id).Scan(
		&todo.ID,
		&todo.Content,
		&todo.UserID,
		&todo.CreatedAt,
	)
	return todo, err
}

// Todoの取得 (Read) (全件取得)
func GetTodos() (todos []Todo, err error) {
	cmd := `SELECT id, content, user_id, created_at FROM todos`

	rows, err := Db.Query(cmd)
	if err != nil {
		log.Fatalln(err)
	}

	for rows.Next() {
		var todo Todo
		err = rows.Scan(
			&todo.ID,
			&todo.Content,
			&todo.UserID,
			&todo.CreatedAt,
		)
		if err != nil {
			log.Fatalln(err)
		}
		todos = append(todos, todo)
	}
	rows.Close()

	return todos, err
}

// Todoの取得 (Read) (特定のユーザーのTodoを取得)
func (u *User) GetTodosByUser() (todos []Todo, err error) {
	cmd := `SELECT id, content, user_id, created_at FROM todos WHERE user_id = $1`

	rows, err := Db.Query(cmd, u.ID)
	if err != nil {
		log.Fatalln(err)
	}

	for rows.Next() {
		var todo Todo
		err = rows.Scan(
			&todo.ID,
			&todo.Content,
			&todo.UserID,
			&todo.CreatedAt,
		)
		if err != nil {
			log.Fatalln(err)
		}
		todos = append(todos, todo)
	}
	rows.Close()

	return todos, err
}

// Todoの更新 (Update)
func (t *Todo) UpdateTodo() error {
	cmd := `UPDATE todos SET content = $1, user_id = $2 WHERE id = $3`
	_, err := Db.Exec(cmd, t.Content, t.UserID, t.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

// Todoの削除 (Update)
func (t *Todo) DeleteTodo() error {
	cmd := `DELETE FROM todos WHERE id = $1`
	_, err := Db.Exec(cmd, t.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}
