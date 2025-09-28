package http

// import (
// 	"database/sql/driver"
// 	"fmt"

// 	"github.com/jackc/pgx/v5/pgtype"
// )

// type GenderT string

// const (
// 	GenderTMale    GenderT = "male"
// 	GenderTFemale  GenderT = "female"
// 	GenderTOther   GenderT = "other"
// 	GenderTUnknown GenderT = "unknown"
// )

// func (e *GenderT) Scan(src interface{}) error {
// 	switch s := src.(type) {
// 	case []byte:
// 		*e = GenderT(s)
// 	case string:
// 		*e = GenderT(s)
// 	default:
// 		return fmt.Errorf("unsupported scan type for GenderT: %T", src)
// 	}
// 	return nil
// }

// type NullGenderT struct {
// 	GenderT GenderT `json:"gender_t"`
// 	Valid   bool    `json:"valid"` // Valid is true if GenderT is not NULL
// }

// // Scan implements the Scanner interface.
// func (ns *NullGenderT) Scan(value interface{}) error {
// 	if value == nil {
// 		ns.GenderT, ns.Valid = "", false
// 		return nil
// 	}
// 	ns.Valid = true
// 	return ns.GenderT.Scan(value)
// }

// // Value implements the driver Valuer interface.
// func (ns NullGenderT) Value() (driver.Value, error) {
// 	if !ns.Valid {
// 		return nil, nil
// 	}
// 	return string(ns.GenderT), nil
// }

// type Address struct {
// 	ID      pgtype.UUID `json:"id"`
// 	Country interface{} `json:"country"`
// 	City    interface{} `json:"city"`
// 	Street  interface{} `json:"street"`
// }

// type Client struct {
// 	ID               pgtype.UUID        `json:"id"`
// 	ClientName       interface{}        `json:"client_name"`
// 	ClientSurname    interface{}        `json:"client_surname"`
// 	Birthday         pgtype.Date        `json:"birthday"`
// 	Gender           GenderT            `json:"gender"`
// 	RegistrationDate pgtype.Timestamptz `json:"registration_date"`
// 	AddressID        pgtype.UUID        `json:"address_id"`
// }

// type Image struct {
// 	ID    pgtype.UUID `json:"id"`
// 	Image []byte      `json:"image"`
// }

// type Product struct {
// 	ID             pgtype.UUID        `json:"id"`
// 	Name           interface{}        `json:"name"`
// 	Category       interface{}        `json:"category"`
// 	Price          pgtype.Numeric     `json:"price"`
// 	AvailableStock int32              `json:"available_stock"`
// 	LastUpdateDate pgtype.Timestamptz `json:"last_update_date"`
// 	SupplierID     pgtype.UUID        `json:"supplier_id"`
// 	ImageID        pgtype.UUID        `json:"image_id"`
// }

// type Supplier struct {
// 	ID          pgtype.UUID `json:"id"`
// 	Name        interface{} `json:"name"`
// 	AddressID   pgtype.UUID `json:"address_id"`
// 	PhoneNumber string      `json:"phone_number"`
// }
