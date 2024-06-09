package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"todo_list/handlers"
	"todo_list/repository"

	_ "github.com/lib/pq"
)

func setupTestServer() *http.ServeMux {
	database, err := SetupTestDB()
	if err != nil {
		log.Fatal(err)
	}
	taskRepo := &repository.TaskRepository{DB: database}
	taskHandler := &handlers.TaskHandler{Repo: taskRepo}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/create", taskHandler.CreateTask)
	mux.HandleFunc("/api/task/", taskHandler.GetTask)
	mux.HandleFunc("/api/tasks/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskHandler.ListTasks(w, r) // Обработка GET запросов
		case http.MethodPut:
			taskHandler.UpdateTask(w, r) // Обработка PUT запросов
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/delete/", taskHandler.DeleteTaskByID)

	return mux
}

func TestCreateTask(t *testing.T) {
	server := setupTestServer()

	task := map[string]interface{}{
		"title":       "Test Task",
		"description": "This is a test task",
		"date":        "2024-05-05",
		"done":        false,
	}

	jsonTask, _ := json.Marshal(task)
	req, err := http.NewRequest("POST", "/api/create", bytes.NewBuffer(jsonTask))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestGetTask(t *testing.T) {
	server := setupTestServer()
	req, err := http.NewRequest("GET", "/api/task/1", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)
	fmt.Println(rr)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

}
func TestUpdateTask(t *testing.T) {
	server := setupTestServer()
	taskUpdate := map[string]interface{}{
		"id":          1,
		"title":       "Test Task Update",
		"description": "This is a test task Update",
		"date":        "2024-08-08",
		"done":        true,
	}
	jsonTask, _ := json.Marshal(taskUpdate)

	req, err := http.NewRequest("PUT", "/api/tasks/", bytes.NewBuffer(jsonTask))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)
	fmt.Println(rr)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestListTasks(t *testing.T) {
	server := setupTestServer()
	for i := 0; i < 50; i++ {
		task := map[string]interface{}{
			"title":       "Test Task Update",
			"description": "This is a test task Update",
			"date":        "2024-08-08",
			"done":        true,
		}
		jsonTask, _ := json.Marshal(task)

		req, err := http.NewRequest("POST", "/api/create", bytes.NewBuffer(jsonTask))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		// Выполняем запрос
		rr := httptest.NewRecorder()
		server.ServeHTTP(rr, req)
	}

	listParam := map[string]interface{}{
		"status": "true",
		"date":   "2024-08-08",
		"limit":  0,
		"offset": 0,
	}

	jsonListParam, _ := json.Marshal(listParam)

	req, err := http.NewRequest("GET", "/api/tasks/", bytes.NewBuffer(jsonListParam))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)
	fmt.Println(rr)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestDeleteTaskByID(t *testing.T) {
	server := setupTestServer()
	req, err := http.NewRequest("DELETE", "/api/delete/10", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)
	fmt.Println(rr)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

}

// Additional test cases for DeleteTaskByID, ListTasks, UpdateTask, etc.
func ConnectTestDB() (*sql.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", dbHost, dbPort, dbUser, dbName, dbPassword)
	return sql.Open("postgres", connStr)
}

func SetupTestDB() (*sql.DB, error) {
	database, err := ConnectTestDB()
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	// Создание тестовой таблицы
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS tasks_test (
        id SERIAL PRIMARY KEY,
        title TEXT,
        description TEXT,
        date DATE,
        done BOOLEAN
    )`
	_, err = database.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	return database, nil
}

func TearDownTestDB(db *sql.DB) {
	// Удаление тестовой таблицы после тестов
	dropTableQuery := `DROP TABLE IF EXISTS tasks_test`
	_, err := db.Exec(dropTableQuery)
	if err != nil {
		fmt.Println("ошибка подключения к БД")
	}
	db.Close()
}
