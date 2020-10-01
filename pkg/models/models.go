package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

type User struct {
	ID             uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; default:uuid_generate_v4()"`
	FirstName      string
	LastName       string
	EmailAddress   string
	HashedPassword []byte
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
	EmailUpdates   bool
	AdvancedView   bool
	Active         bool
}

type PasswordResetToken struct {
	ID           uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; default:uuid_generate_v4()"`
	EmailAddress string
	CreatedAt    time.Time
}

type Note struct {
	Content       string
	Category      string
	Tags          []string
	DueDate       time.Time
	RemindDate    time.Time
	Completed     bool
	CompletedDate time.Time
	//how to connect to user and project
	//how to add a dependent on schema
}
