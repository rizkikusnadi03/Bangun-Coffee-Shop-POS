package entity

import "time"

// Role constants
const (
	RoleOwner   = "owner"
	RoleCashier = "cashier"
)

// User merepresentasikan entitas pengguna dalam tabel users.
type User struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	Role      string     `json:"role"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// IsOwner memeriksa apakah user memiliki role owner.
func (u *User) IsOwner() bool {
	return u.Role == RoleOwner
}

// IsCashier memeriksa apakah user memiliki role cashier.
func (u *User) IsCashier() bool {
	return u.Role == RoleCashier
}

// ToResponse mengonversi entitas User ke format UserResponse yang aman dikirim ke client.
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
	}
}

// LoginRequest merepresentasikan payload permintaan login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// RegisterRequest merepresentasikan payload permintaan pendaftaran/pembuatan user kasir.
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=72"`
	Role     string `json:"role" validate:"required,oneof=owner cashier"`
}

// UserResponse merepresentasikan format data user yang dikembalikan ke client tanpa field sensitif.
type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// LoginResponse merepresentasikan payload respon sukses login beserta token JWT.
type LoginResponse struct {
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expires_at"`
	User      UserResponse `json:"user"`
}
