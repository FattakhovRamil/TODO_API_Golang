package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", dbHost, dbPort, dbUser, dbName, dbPassword)
	database, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}
	err = database.Ping()
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	createTableQuery := `
    CREATE TABLE IF NOT EXISTS tasks (
        id SERIAL PRIMARY KEY,
        title TEXT,
        description TEXT,
        date DATE,
        done BOOLEAN
    )`
	_, err = database.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания таблицы tasks: %w", err)
	}

	fmt.Println("Подключение к БД прошло успешно. Таблица создана")
	return database, nil
}
