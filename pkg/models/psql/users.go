package psql

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/cpustejovsky/estuary/pkg/models"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(first_name, last_name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (first_name, last_name, email, hashed_password) VALUES($1, $2, $3, $4)`

	_, err = m.DB.Exec(stmt, first_name, last_name, email, string(hashedPassword))
	if err != nil {
		var postgresError *pq.Error
		if errors.As(err, &postgresError) {
			if postgresError.Code == "23505" {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (string, error) {
	var id uuid.UUID
	var hashedPassword []byte
	stmt := "SELECT id, hashed_password FROM users WHERE email = $1 AND active = TRUE"
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "error", models.ErrInvalidCredentials
		} else {
			return "error", err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "error", models.ErrInvalidCredentials
		} else {
			return "error", err
		}
	}

	return id.String(), nil
}

func (m *UserModel) Get(id string) (*models.User, error) {
	u := &models.User{}
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	stmt := `SELECT id, first_name, last_name, email, created, active FROM users WHERE id = $1`
	err = m.DB.QueryRow(stmt, uuid).Scan(&u.ID, &u.FirstName, &u.LastName, &u.EmailAddress, &u.CreatedAt, &u.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return u, nil
}
