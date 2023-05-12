package main

import "fmt"

func NewDatabaseError(text string, err error) error {
	if err != nil {
		return fmt.Errorf("Database error: %s - %w", DatabaseError{Message: text}, err)
	}
	return &DatabaseError{Message: text}
}

type DatabaseError struct {
	Message string
}

func (e *DatabaseError) Error() string {
	return e.Message
}
