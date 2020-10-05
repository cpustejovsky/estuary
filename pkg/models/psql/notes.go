package psql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/cpustejovsky/estuary/pkg/models"
)

type NotesModel struct {
	DB *sql.DB
}

func (m *NotesModel) Insert(content, accountId string) error {
	uuid, err := uuid.Parse(accountId)
	if err != nil {
		return err
	}
	stmt := `
	INSERT INTO notes (content, account_id) 
	VALUES($1, $2)`

	_, err = m.DB.Exec(stmt, content, uuid)
	if err != nil {
		return err
	}

	return nil
}

func (m *NotesModel) GetByCategory(accountId, category string) (*[]models.Note, error) {
	var notes []models.Note
	uuid, err := uuid.Parse(accountId)
	if err != nil {
		return nil, err
	}
	stmt := `
	SELECT content, category, tags, created, due_date
	FROM notes 
	WHERE account_id = $1 AND category = $2`
	rows, err := m.DB.Query(stmt, uuid, category)
	defer rows.Close()
	for rows.Next() {
		n := &models.Note{}
		err = rows.Scan(&n.Content, &n.Category, &n.Tags, &n.Created, &n.DueDate)
		notes = append(notes, *n)
	}
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return &notes, nil
}
