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

func (m *NoteModel) Insert(accountId, content string) (*models.Note, error) {
	accountUUID, err := uuid.Parse(accountId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	query := `
	INSERT INTO notes (account_id, content) 
	VALUES($1, $2) RETURNING note_id, content, category, created`
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer stmt.Close()
	n := &models.Note{}
	err = stmt.QueryRow(accountUUID, content).Scan(&n.ID, &n.Content, &n.Category, &n.Created)
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

func (m *NoteModel) Update(accountId, noteId, content string) (*models.Note, error) {
	accountUUID, err := uuid.Parse(accountId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	noteUUID, err := uuid.Parse(noteId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	query := `
	UPDATE notes
	SET content = $3
	WHERE account_id = $1 AND note_id = $2 RETURNING note_id, content, category, created
	`
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer stmt.Close()
	n := &models.Note{}
	err = stmt.QueryRow(accountUUID, noteUUID, content).Scan(&n.ID, &n.Content, &n.Category, &n.Created)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return n, nil

}

func (m *NoteModel) Delete(noteId string) error {
	sqlStatement := `
	DELETE FROM notes
	WHERE note_id = $1;`
	_, err := m.DB.Exec(sqlStatement, noteId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrNoRecord
		} else {
			return err
		}
	}
	return nil
}
