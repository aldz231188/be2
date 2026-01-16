package domain

import (
	"time"

	"github.com/google/uuid"
)

type Gender int

const (
	MALE Gender = iota
	FEMALE
)

type Client struct {
	UserID        uuid.UUID
	ID            uuid.UUID
	ClientName    string
	ClientSurname string
	// Birthday         time.Time
	// Gender           Gender
	// RegistrationDate time.Time
	// Address          uuid.UUID
}

type Address struct {
	ID      uuid.UUID
	Country string
	City    string
	Street  string
}

type User struct {
	ID           uuid.UUID
	Username     string
	PasswordHash string
	CreatedAt    time.Time
	TokenVersion int32
}

type Session struct {
	JTIHash   string
	UserID    uuid.UUID
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}
