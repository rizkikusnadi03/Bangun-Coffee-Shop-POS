package entity

import "time"

// Category merepresentasikan entitas kategori produk dalam tabel categories.
type Category struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

// CreateCategoryRequest merepresentasikan payload permintaan pembuatan kategori baru.
type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}

// UpdateCategoryRequest merepresentasikan payload permintaan pembaruan data kategori.
type UpdateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}
