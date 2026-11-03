package entity

import "time"

// Konstanta tipe pergerakan stok
const (
	MovementTypeIn         = "in"
	MovementTypeOut        = "out"
	MovementTypeAdjustment = "adjustment"
)

// StockMovement merepresentasikan audit trail pergerakan stok dalam tabel stock_movements.
type StockMovement struct {
	ID          string    `json:"id"`
	ProductID   string    `json:"product_id"`
	Product     *Product  `json:"product,omitempty"`
	UserID      string    `json:"user_id"`
	User        *User     `json:"user,omitempty"`
	Type        string    `json:"type"`
	Quantity    int       `json:"quantity"`
	StockBefore int       `json:"stock_before"`
	StockAfter  int       `json:"stock_after"`
	Notes       string    `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// StockAdjustmentRequest merepresentasikan payload permintaan penyesuaian stok produk.
type StockAdjustmentRequest struct {
	Quantity int    `json:"quantity" validate:"required,min=1"`
	Type     string `json:"type" validate:"required,oneof=in out adjustment"`
	Notes    string `json:"notes" validate:"max=500"`
}
