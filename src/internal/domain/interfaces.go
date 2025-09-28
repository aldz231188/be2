package domain

import "context"

type Service interface {
	AddAddress(ctx context.Context, c *Adress) error
	// SaveAdress(c *Adress) error
}

type Repo interface {
	AddAddress(ctx context.Context, c *Adress) error
	// SaveAdress(c *Adress) error
}

// type AdressRepo interface {

// }
