package mock

import (
	"time"

	"github.com/cpustejovsky/estuary/pkg/models"
	"github.com/google/uuid"
)

var mockAccount = &models.Account{
	ID:           uuid.New(),
	FirstName:    "Charles",
	LastName:     "Pustejovsky",
	EmailAddress: "charles.pustejovsky@gmail.com",
	CreatedAt:    time.Now(),
	EmailUpdates: true,
	AdvancedView: true,
}
