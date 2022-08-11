package user

import (
	"context"

	"github.com/1995parham-teaching/fandogh/internal/model"
)

type MemoryUser struct {
	store map[string]model.User
}

func NewMemoryUser() *MemoryUser {
	return &MemoryUser{
		store: make(map[string]model.User),
	}
}

func (m MemoryUser) Set(ctx context.Context, user model.User) error {
	if _, ok := m.store[user.Email]; ok {
		return ErrEmailDuplicate
	}

	m.store[user.Email] = user

	return nil
}

func (m MemoryUser) Get(ctx context.Context, email string) (model.User, error) {
	user, ok := m.store[email]
	if ok {
		return user, nil
	}

	return user, ErrEmailNotFound
}
