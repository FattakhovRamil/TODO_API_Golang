package models_test

import (
	"testing"
	. "todo_list/models"
)

func TestValidate(t *testing.T) {
	task := Task{
		ID:          1,
		Title:       "", // Пустой заголовок
		Description: "Description",
		Date:        "2022-06-08",
		Done:        false,
	}

	// Проверяем, что Validate возвращает ошибку, так как заголовок пустой
	err := task.Validate()
	if err == nil {
		t.Error("Ожидалась ошибка, так как заголовок пустой, но получено nil")
	}

	task = Task{
		ID:          1,
		Title:       "заголовок",
		Description: "", // Пустое описание
		Date:        "2022-06-08",
		Done:        false,
	}

	// Проверяем, что Validate возвращает ошибку, так как описание пустое
	err = task.Validate()
	if err == nil {
		t.Error("Ожидалась ошибка, так как описание пустое, но получено nil")
	}

	task = Task{
		ID:          1,
		Title:       "заголовок",
		Description: "описание",
		Date:        "202-06-08", // Неправильная дата
		Done:        false,
	}

	// Проверяем, что Validate возвращает ошибку, так как описание пустое
	err = task.Validate()
	if err == nil {
		t.Error("Ожидалась ошибка, так как дата не корректна, но получено nil")
	}

	task = Task{
		ID:          1,
		Title:       "заголовок",
		Description: "описание",
		Date:        "", // Пустая дата
		Done:        false,
	}

	// Проверяем, что Validate возвращает ошибку, так как описание пустое
	err = task.Validate()
	if err == nil {
		t.Error("Ожидалась ошибка, так как дата пуста, но получено nil")
	}

	// Создаем экземпляр задачи с правильными данными
	task = Task{
		ID:          1,
		Title:       "Title",
		Description: "Description",
		Date:        "2022-06-08",
		Done:        false,
	}

	// Проверяем, что Validate не возвращает ошибку для правильных данных
	err = task.Validate()
	if err != nil {
		t.Errorf("Не ожидалась ошибка, но получено: %v", err)
	}
}

func TestSetDate(t *testing.T) {
	var dateT string = "2024-05-05"
	task := Task{
		ID:          1,
		Title:       "Title",
		Description: "Description",
		Date:        "2022-06-08",
		Done:        false,
	}

	err := task.SetDate(dateT)

	if err != nil {
		t.Errorf("Не ожидалась ошибка, но получено: %v", err)
	}

	dateT = "2024-05ук-05"

	err = task.SetDate(dateT)

	if err == nil {
		t.Error("Ожидалась ошибка, так как строка неверного формата но получено nil")
	}

}

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







