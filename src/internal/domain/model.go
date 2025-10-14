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
	ID               uuid.UUID
	ClientName       string
	ClientSurname    string
	Birthday         time.Time
	Gender           Gender
	RegistrationDate time.Time
	Address          uuid.UUID
}

type Address struct {
	ID      uuid.UUID
	Country string
	City    string
	Street  string
}
