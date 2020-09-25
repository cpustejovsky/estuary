package psql

import (
	"database/sql"
	"errors"

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

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	stmt := "SELECT id, hashed_password FROM users WHERE email = $2 AND active = TRUE"
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	u := &models.User{}

	stmt := `SELECT id, first_name, last_name, email, created, active FROM users WHERE id = $1`
	err := m.DB.QueryRow(stmt, id).Scan(&u.ID, &u.FirstName, &u.LastName, &u.EmailAddress, &u.CreatedAt, &u.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return u, nil
}
