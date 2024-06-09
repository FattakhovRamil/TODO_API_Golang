package repository

import (
	"database/sql"

	"todo_list/models"
)

type TaskRepository struct {
	DB *sql.DB
}

func (r *TaskRepository) Create(task *models.Task) error {
	query := "INSERT INTO tasks (title, description, date, done) VALUES ($1, $2, $3, $4) RETURNING id"
	return r.DB.QueryRow(query, task.Title, task.Description, task.Date, task.Done).Scan(&task.ID)
}

func (r *TaskRepository) Get(id int) (*models.Task, error) {
	task := &models.Task{}
	query := "SELECT * FROM tasks WHERE id=$1"
	err := r.DB.QueryRow(query, id).Scan(&task.ID, &task.Title, &task.Description, &task.Date, &task.Done)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return task, err
}

func (r *TaskRepository) Update(task *models.Task) error {
	query := "UPDATE tasks SET title=$1, description=$2, date=$3, done=$4 WHERE id=$5"
	_, err := r.DB.Exec(query, task.Title, task.Description, task.Date, task.Done, task.ID)
	return err
}

func (r *TaskRepository) Delete(id int) error {
	query := "DELETE FROM tasks WHERE id=$1"
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *TaskRepository) List(status string, date string, limit int, offset int) ([]*models.Task, error) {
	tasks := []*models.Task{}
	query := `SELECT id, title, description, date, done FROM tasks WHERE done = $1 AND date = $2 LIMIT $3 OFFSET $4`
	rows, err := r.DB.Query(query, status, date, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		task := &models.Task{}
		err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.Date, &task.Done)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
