package psql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/cpustejovsky/estuary/pkg/models"
)

type ResetTokenModel struct {
	DB *sql.DB
}

func (m *ResetTokenModel) Insert(id uuid.UUID, email string) error {
	stmt := `
	INSERT INTO reset_tokens (id, email) 
	VALUES($1, $2)`

	_, err := m.DB.Exec(stmt, id, email)
	if err != nil {
		return err
	}

	return nil
}

func (m *ResetTokenModel) Get(id, email string) (*models.ResetToken, error) {
	r := &models.ResetToken{}
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	stmt := `
	SELECT id, email, created 
	FROM reset_tokens 
	WHERE id = $1 AND email = $2`
	err = m.DB.QueryRow(stmt, uuid, email).Scan(&r.ID, &r.EmailAddress, &r.CreatedAt)
	hourAgo := time.Now().Add(-1 * time.Hour)
	if hourAgo.After(r.CreatedAt) {
		return nil, errors.New("expired token")
	}
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return r, nil
}
