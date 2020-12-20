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

type Account struct {
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

type ResetToken struct {
	ID           uuid.UUID
	EmailAddress string
	CreatedAt    time.Time
}

type Note struct {
	ID            uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; default:uuid_generate_v4()"`
	Content       string
	Category      string
	Tags          []string
	Created       time.Time
	DueDate       time.Time
	RemindDate    time.Time
	Completed     bool
	CompletedDate time.Time
	//Connect to Account
	AccountID uuid.UUID
	//Connect to Dependents
}
