package common

import (
	"fmt"
	"strings"
)

func NewDatabaseError(text string, err error) error {
	if err != nil {
		return fmt.Errorf("database error: %s - %w", DatabaseError{Message: text}, err)
	}
	return &DatabaseError{Message: text}
}

type DatabaseError struct {
	Message string
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("database error: %s", e.Message)
}

type ConfigLoadError struct {
	err   error
	label string
}

func (e *ConfigLoadError) Error() string {
	builder := strings.Builder{}

	builder.WriteString("failed to load config label ")
	builder.WriteString(e.label)
	if e.err != nil {
		builder.WriteString(": ")
		builder.WriteString(e.err.Error())
	}

	return builder.String()
}

func NewConfigLoadError(label string, err error) error {
	return &ConfigLoadError{label: label, err: err}
}
