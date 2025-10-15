// Package domain contains auth service repository interfaces
package domain

import (
	"github.com/google/uuid"
)

// UserRepository interface for user data access
type UserRepository interface {
	FindByID(userID uuid.UUID) (*User, error)
	FindByZaloID(zaloID string) (*User, error)
	FindByEmail(email string) (*User, error)
	FindByPhone(phone string) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(userID uuid.UUID) error
}