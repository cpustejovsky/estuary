package psql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/cpustejovsky/estuary/pkg/models"
)

type NoteModel struct {
	DB *sql.DB
}

func (m *NoteModel) Insert(accountId, content string) (string, error) {
	accountUUID, err := uuid.Parse(accountId)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	query := `
	INSERT INTO notes (account_id, content) 
	VALUES($1, $2) RETURNING note_id`
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer stmt.Close()
	var noteId uuid.UUID
	err = stmt.QueryRow(accountUUID, content).Scan(&noteId)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return noteId.String(), nil
}

func (m *NoteModel) Get(noteId string) (*models.Note, error) {
	uuid, err := uuid.Parse(noteId)
	if err != nil {
		return nil, err
	}
	stmt := `
	SELECT note_id, content, category, created
	FROM notes 
	WHERE note_id = $1`
	n := &models.Note{}
	row := m.DB.QueryRow(stmt, uuid)
	err = row.Scan(&n.ID, &n.Content, &n.Category, &n.Created)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return n, nil
}

func (m *NoteModel) GetByCategory(accountId, category string) (*[]models.Note, error) {
	var notes []models.Note
	uuid, err := uuid.Parse(accountId)
	if err != nil {
		return nil, err
	}
	stmt := `
	SELECT note_id, content, category, created
	FROM notes 
	WHERE account_id = $1 AND category = $2`
	rows, err := m.DB.Query(stmt, uuid, category)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		n := &models.Note{}
		err = rows.Scan(&n.ID, &n.Content, &n.Category, &n.Created)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
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
