package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	_ "todo_list/docs"
	"todo_list/models"
	"todo_list/repository"
)

// TaskHandler represents the handler for task-related operations.
type TaskHandler struct {
	Repo *repository.TaskRepository // Repo is a repository for task-related operations.
}

// CreateTask обрабатывает создание новой задачи.
// @Summary Создание новой задачи
// @Description Создание новой задачи с предоставленными данными
// @Accept json
// @Produce json
// @Param task body Task true "Объект задачи для создания"
// @Success 201 {object} Task
// @Failure 400 {object} error "Неверный формат запроса"
// @Failure 500 {object} error "Внутренняя ошибка сервера"
// @Router /api/create [post]
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var task models.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := task.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Repo.Create(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GetTask обрабатывает получение задачи по ее ID.
// @Summary Получение задачи по ID
// @Description Получение задачи по ее ID
// @Accept json
// @Produce json
// @Param id path int true "Идентификатор задачи"
// @Success 200 {object} Task
// @Failure 400 {object} error "Неверный формат запроса"
// @Failure 404 {object} error "Задача не найдена"
// @Failure 500 {object} error "Внутренняя ошибка сервера"
// @Router /api/task/{id} [get]
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Path[len("/api/task/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	task, err := h.Repo.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	json.NewEncoder(w).Encode(task)
}

// UpdateTask обрабатывает обновление существующей задачи.
// @Summary Обновление задачи
// @Description Обновление существующей задачи
// @Accept json
// @Produce json
// @Param task body Task true "Объект задачи для обновления"
// @Success 200 {object} Task
// @Failure 400 {object} error "Неверный формат запроса"
// @Failure 500 {object} error "Внутренняя ошибка сервера"
// @Router /api/tasks [put]
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := task.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.Repo.Update(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteTaskByID обрабатывает удаление задачи по ее ID.
// @Summary Удаление задачи
// @Description Удаление задачи по ее ID
// @Accept json
// @Produce json
// @Param id path int true "Идентификатор задачи"
// @Success 200 {string} string "Задача успешно удалена"
// @Failure 400 {object} error "Неверный ID"
// @Failure 404 {object} error "Задача не найдена"
// @Router /api/delete/{id} [delete]
func (h *TaskHandler) DeleteTaskByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.URL.Path[len("/api/delete/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = h.Repo.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

// ListTasks обрабатывает список задач на основе определенных критериев.
// @Summary Получение списка задач
// @Description Получение списка задач на основе определенных критериев
// @Accept json
// @Produce json
// @Param status query string true "Статус задачи"
// @Param date query string true "Дата"
// @Param limit query int true "Лимит"
// @Param offset query int true "Смещение"
// @Success 200 {object} []Task "Список задач"
// @Failure 400 {object} error "Неверный формат даты или параметров"
// @Failure 500 {object} error "Внутренняя ошибка сервера"
// @Router /api/tasks [get]
func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var listReq struct {
		Status string `json:"status"`
		Date   string `json:"date"`
		Limit  int    `json:"limit"`
		Offset int    `json:"offset"`
	}

	err := json.NewDecoder(r.Body).Decode(&listReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if listReq.Status == "" || listReq.Date == "" {
		http.Error(w, "Status and date are required", http.StatusBadRequest)
		return
	}

	_, err = time.Parse("2006-01-02", listReq.Date)
	if err != nil {
		http.Error(w, "Invalid date format, expected YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	if listReq.Limit < 0 || listReq.Offset < 0 {
		http.Error(w, "Limit and offset must be non-negative", http.StatusBadRequest)
		return
	}

	tasks, err := h.Repo.List(listReq.Status, listReq.Date, listReq.Limit, listReq.Offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tasks)
}

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Done        bool   `json:"done"`
}

// Validate checks if the task data is valid.
// func (t *Task) Validate() error {
// 	// Check if the title is empty.
// 	if t.Title == "" {
// 		return errors.New("title is required")
// 	}
// 	// Check if the description is empty.
// 	if t.Description == "" {
// 		return errors.New("description is required")
// 	}
// 	// Check if the date is empty.
// 	if t.Date == "" {
// 		return errors.New("valid date is required")
// 	}

// 	// Parse the date and check if it's in the correct format (YYYY-MM-DD).
// 	_, err := time.Parse("2006-01-02", t.Date)
// 	if err != nil {
// 		return errors.New("invalid date format, expected YYYY-MM-DD")
// 	}
// 	return nil
// }

// // SetDate sets the date of the task.
// func (t *Task) SetDate(dateString string) error {
// 	// Parse the date string.
// 	date, err := time.Parse("2006-01-02", dateString)
// 	if err != nil {
// 		return errors.New("invalid date format, expected YYYY-MM-DD")
// 	}
// 	// Format the date as YYYY-MM-DD and set it to the task.
// 	t.Date = date.Format("2006-01-02")
// 	return nil
// }


