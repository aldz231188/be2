package domain

import (
	"time"

	"github.com/google/uuid"
)

type Gender int

const (
	UNKNOWN Gender = iota
	MALE
	FEMALE
)

type Client struct {
	Id               uuid.UUID
	ClientName       string
	ClientSurname    string
	Birthday         time.Time
	Gender           Gender
	RegistrationDate time.Time
	Address          Address
}

type Address struct {
	Id      uuid.UUID
	Country string
	City    string
	Street  string
}
