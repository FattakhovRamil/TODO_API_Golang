package server

import (
	"fmt"
	"log"
	"net/http"
	db "todo_list/db"
	_ "todo_list/docs"
	"todo_list/handlers"
	"todo_list/repository"

	httpSwagger "github.com/swaggo/http-swagger"
)

func ServerStart() error {
	database, err := db.ConnectDB()

	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		fmt.Println("Error opening database:", err)
		return err
	}
	defer database.Close()
	taskRepo := &repository.TaskRepository{DB: database}
	taskHandler := &handlers.TaskHandler{Repo: taskRepo}

	http.HandleFunc("/api/create", taskHandler.CreateTask)      // Create POST
	http.HandleFunc("/api/task/", taskHandler.GetTask)          // Get POST
	http.HandleFunc("/api/delete/", taskHandler.DeleteTaskByID) // Delete DELETE
	http.HandleFunc("/api/tasks/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskHandler.ListTasks(w, r) // Обработка GET запросов
		case http.MethodPut:
			taskHandler.UpdateTask(w, r) // Обработка PUT запросов
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	fmt.Println("Server start listening on port 3001")
	err = http.ListenAndServe(":3001", nil)
	fmt.Println("Server end listening on port 3001")
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
