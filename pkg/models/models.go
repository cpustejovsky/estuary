package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

type Note struct {
	Content  string
	Category string
	Tags     []string
	DueDate  time.Time
	RemindDate time.Time
	Completed bool
	CompletedDate time.Time
	//how to connect to user and project
	//how to add a dependent on schema
}
