package mocks

import (
	"github.com/Jackob2004/snippetbox/internal/models"
)

type UserModel struct{}

const Email = "alice@example.com"
const Password = "pa$$word"

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case Email:
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == Email && password == Password {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}
