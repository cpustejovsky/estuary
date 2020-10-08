package mock

import (
	"time"

	"github.com/cpustejovsky/estuary/pkg/models"
	"github.com/google/uuid"
)

var acctID = uuid.New()

var mockAccount = &models.Account{
	ID:           acctID,
	FirstName:    "Charles",
	LastName:     "Pustejovsky",
	EmailAddress: "charles.pustejovsky@gmail.com",
	CreatedAt:    time.Now(),
	EmailUpdates: true,
	AdvancedView: true,
}

type AccountModel struct{}

func (m *AccountModel) Insert(first_name, last_name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *AccountModel) Authenticate(email, password string) (int, error) {
	switch email {
	case "charles.pustejovsky@gmail.com":
		return 1, nil
	default:
		return 0, models.ErrInvalidCredentials
	}
}

func (m *AccountModel) Get(id uuid.UUID) (*models.Account, error) {
	switch id {
	case acctID:
		return mockAccount, nil
	default:
		return nil, models.ErrNoRecord
	}
}
