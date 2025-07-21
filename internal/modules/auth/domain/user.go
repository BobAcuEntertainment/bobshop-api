package domain

import (
	"time"

	"github.com/google/uuid"

	"bobshop/pkg/auth"
)

type User struct {
	ID           uuid.UUID `bson:"_id"`
	Email        string    `bson:"email"`
	PasswordHash string    `bson:"password_hash"`
	Role         Role      `bson:"role"`
	CreatedAt    time.Time `bson:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at"`
}

func NewUser(email, password string) (*User, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	user := &User{
		ID:        id,
		Email:     email,
		Role:      UserRole,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return err
	}
	u.PasswordHash = hashedPassword
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) CheckPassword(password string) bool {
	return auth.ComparePassword(u.PasswordHash, password)
}
