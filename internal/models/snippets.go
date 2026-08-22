package models

import (
	"database/sql"
	"errors"
	"time"
)

type SnippetModelInterface interface {
	Insert(title string, content string, expires, userId int) (int, error)
	Get(id int) (Snippet, error)
	Delete(snippetId, userId int) error
	Update(snippetId, userId int, title string, content string, expires int) error
	Latest() ([]Snippet, error)
	GetSnippets(userId int) ([]Snippet, error)
}

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
	Creator SnippetCreator
}

type SnippetCreator struct {
	UserID   int
	UserName string
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Update(snippetId, userId int, title string, content string, expires int) error {
	stmt := `UPDATE snippets SET title = ?, content = ?, expires = DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY)
    WHERE id = ? AND user_id = ?`

	res, err := m.DB.Exec(stmt, title, content, expires, snippetId, userId)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNoRecord
	}

	return nil
}

func (m *SnippetModel) Delete(snippetId, userId int) error {
	stmt := `DELETE FROM snippets WHERE id = ? AND user_id = ?`

	res, err := m.DB.Exec(stmt, snippetId, userId)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNoRecord
	}

	return nil
}

func (m *SnippetModel) Insert(title, content string, expires, userId int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires, user_id)
	VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY), ?)`

	result, err := m.DB.Exec(stmt, title, content, expires, userId)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (Snippet, error) {
	stmt := `SELECT snippets.id, title, content, snippets.created, expires, user_id, name FROM snippets
    JOIN users ON snippets.user_id = users.id
    WHERE expires > UTC_TIMESTAMP() AND snippets.id = ?`

	row := m.DB.QueryRow(stmt, id)
	var s Snippet

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires, &s.Creator.UserID, &s.Creator.UserName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		}

		return Snippet{}, err
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]Snippet, error) {
	stmt := `SELECT snippets.id, title, content, snippets.created, expires, user_id, name FROM snippets
    JOIN users ON snippets.user_id = users.id
	WHERE expires > UTC_TIMESTAMP() ORDER BY snippets.id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return []Snippet{}, err
	}
	defer rows.Close()

	return mapRows(rows)
}

func (m *SnippetModel) GetSnippets(userId int) ([]Snippet, error) {
	stmt := `SELECT snippets.id, title, content, snippets.created, expires, user_id, name FROM snippets
    JOIN users ON snippets.user_id = users.id 
    WHERE user_id = ?`

	rows, err := m.DB.Query(stmt, userId)
	if err != nil {
		return []Snippet{}, err
	}
	defer rows.Close()

	return mapRows(rows)
}

func mapRows(rows *sql.Rows) ([]Snippet, error) {
	var snippets []Snippet
	for rows.Next() {
		var s Snippet
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires, &s.Creator.UserID, &s.Creator.UserName)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
