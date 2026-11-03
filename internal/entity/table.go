package entity

import "time"

// Konstanta status meja
const (
	TableStatusAvailable = "available"
	TableStatusOccupied  = "occupied"
	TableStatusReserved  = "reserved"
)

// Table merepresentasikan entitas meja restoran dalam tabel tables.
type Table struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Capacity  int        `json:"capacity"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

// CreateTableRequest merepresentasikan payload pembuatan meja baru.
type CreateTableRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=100"`
	Capacity int    `json:"capacity" validate:"required,min=1,max=50"`
}

// UpdateTableRequest merepresentasikan payload pembaruan data meja.
type UpdateTableRequest struct {
	Name     string `json:"name" validate:"omitempty,min=1,max=100"`
	Capacity int    `json:"capacity" validate:"omitempty,min=1,max=50"`
	Status   string `json:"status" validate:"omitempty,oneof=available occupied reserved"`
}
