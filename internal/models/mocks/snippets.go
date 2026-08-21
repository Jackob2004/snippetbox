package mocks

import (
	"database/sql"
	"time"

	"github.com/Jackob2004/snippetbox/internal/models"
)

var mockSnippet = models.Snippet{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnippetModel struct{}

func (m *SnippetModel) Delete(id int) error {
	if id != 1 {
		return models.ErrNoRecord
	}

	return nil
}

func (m *SnippetModel) Insert(title string, content string, expires, userId int) (int, error) {
	return 2, nil
}

func (m *SnippetModel) Get(id int) (models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return models.Snippet{}, models.ErrNoRecord
	}
}

func (m *SnippetModel) Latest() ([]models.Snippet, error) {
	return []models.Snippet{mockSnippet}, nil
}

func (m *SnippetModel) GetSnippets(userId int) ([]models.Snippet, error) {
	if userId == 1 {
		return []models.Snippet{mockSnippet}, nil
	}

	return []models.Snippet{}, sql.ErrNoRows
}
