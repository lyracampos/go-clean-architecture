package entities

import "time"

type User struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
	Role      string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(firstName, lastName, email, role string) *User {
	return &User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Role:      role,
	}
}
