package db_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	db "todo_list/db"

	_ "github.com/lib/pq"
)

func TestConnectDB(t *testing.T) {
	// Подключение к тестовой базе данных.
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", dbHost, dbPort, dbUser, dbName, dbPassword)
	database, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Тест Ошибка при открытии тестовой БД: %v", err)
	}
	defer database.Close()

	// Проверка успешного подключения.
	err = database.Ping()
	if err != nil {
		t.Fatalf("Тест ошибка при проверке соединения: %v", err)
	}

	// Проверка создания таблицы tasks_test.
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
		t.Fatalf("Тест ошибка при создании таблицы tasks_test: %v", err)
	}

	// Вызов тестируемой функции.
	_, err = db.ConnectDB()
	if err != nil {
		t.Fatalf("Тест ошибка при подключении к БД и создании таблицы: %v", err)
	}
	TearDownTestDB(t, database)
}



func TearDownTestDB(t *testing.T, db *sql.DB) {
	// Удаление тестовой таблицы после тестов
	dropTableQuery := `DROP TABLE IF EXISTS tasks_test`
	_, err := db.Exec(dropTableQuery)
	if err != nil {
		t.Fatalf("ошибка при удалении таблицы tasks_test: %v", err)
	}
	db.Close()
}
