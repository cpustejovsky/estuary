package psql

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/cpustejovsky/estuary/pkg/models"
)

type ResetTokenModel struct {
	DB *sql.DB
}

func (m *ResetTokenModel) Insert(email string) error {
	stmt := `
	INSERT INTO password_reset_tokens (email) 
	VALUES($1)`

	_, err := m.DB.Exec(stmt, email)
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
	FROM password_reset_tokens 
	WHERE id = $1 AND email = $2`
	err = m.DB.QueryRow(stmt, uuid, email).Scan(&r.ID, &r.EmailAddress, &r.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return r, nil
}
