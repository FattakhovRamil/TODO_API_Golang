package models

import (
	"errors"
	"time"
	_ "todo_list/docs" // Importing docs for Swagger documentation generation.
)

// @name Task
type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Done        bool   `json:"done"`
}

// Validate checks if the task data is valid.
func (t *Task) Validate() error {
	// Check if the title is empty.
	if t.Title == "" {
		return errors.New("title is required")
	}
	// Check if the description is empty.
	if t.Description == "" {
		return errors.New("description is required")
	}
	// Check if the date is empty.
	if t.Date == "" {
		return errors.New("valid date is required")
	}

	// Parse the date and check if it's in the correct format (YYYY-MM-DD).
	_, err := time.Parse("2006-01-02", t.Date)
	if err != nil {
		return errors.New("invalid date format, expected YYYY-MM-DD")
	}
	return nil
}

// SetDate sets the date of the task.
func (t *Task) SetDate(dateString string) error {
	// Parse the date string.
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return errors.New("invalid date format, expected YYYY-MM-DD")
	}
	// Format the date as YYYY-MM-DD and set it to the task.
	t.Date = date.Format("2006-01-02")
	return nil
}
