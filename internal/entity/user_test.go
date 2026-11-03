package entity

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"web-coffee/pkg/validator"
)

func TestUserMethods(t *testing.T) {
	now := time.Now()
	owner := User{
		ID:        "user-1",
		Name:      "Owner Name",
		Email:     "owner@example.com",
		Password:  "hashed-password",
		Role:      RoleOwner,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if !owner.IsOwner() {
		t.Errorf("expected IsOwner() to be true")
	}
	if owner.IsCashier() {
		t.Errorf("expected IsCashier() to be false for owner")
	}

	cashier := User{
		ID:       "user-2",
		Role:     RoleCashier,
		IsActive: true,
	}

	if !cashier.IsCashier() {
		t.Errorf("expected IsCashier() to be true")
	}
	if cashier.IsOwner() {
		t.Errorf("expected IsOwner() to be false for cashier")
	}

	// Test ToResponse
	resp := owner.ToResponse()
	if resp.ID != owner.ID || resp.Name != owner.Name || resp.Email != owner.Email || resp.Role != owner.Role {
		t.Errorf("mismatch in ToResponse conversion: %+v vs %+v", resp, owner)
	}

	// Verify JSON marshal never includes password
	ownerBytes, err := json.Marshal(owner)
	if err != nil {
		t.Fatalf("failed to marshal User: %v", err)
	}
	if strings.Contains(string(ownerBytes), "password") || strings.Contains(string(ownerBytes), "hashed-password") {
		t.Errorf("expected password to be omitted from JSON, got: %s", string(ownerBytes))
	}
}

func TestUserRequestsValidation(t *testing.T) {
	v := validator.New()

	t.Run("LoginRequest valid", func(t *testing.T) {
		req := LoginRequest{
			Email:    "test@example.com",
			Password: "secretpassword",
		}
		if errs := v.Validate(req); errs != nil {
			t.Errorf("expected valid login request, got errs: %v", errs)
		}
	})

	t.Run("LoginRequest invalid", func(t *testing.T) {
		req := LoginRequest{
			Email:    "invalid-email",
			Password: "123", // too short (< 6)
		}
		errs := v.Validate(req)
		if errs == nil {
			t.Fatal("expected validation errors, got nil")
		}
		if _, ok := errs["email"]; !ok {
			t.Errorf("expected email error, got: %v", errs)
		}
		if _, ok := errs["password"]; !ok {
			t.Errorf("expected password error, got: %v", errs)
		}
	})

	t.Run("RegisterRequest valid", func(t *testing.T) {
		req := RegisterRequest{
			Name:     "Kasir Satu",
			Email:    "kasir1@example.com",
			Password: "password123",
			Role:     RoleCashier,
		}
		if errs := v.Validate(req); errs != nil {
			t.Errorf("expected valid register request, got errs: %v", errs)
		}
	})

	t.Run("RegisterRequest invalid role", func(t *testing.T) {
		req := RegisterRequest{
			Name:     "Super Admin",
			Email:    "super@example.com",
			Password: "password123",
			Role:     "superadmin", // not in owner, cashier
		}
		errs := v.Validate(req)
		if errs == nil {
			t.Fatal("expected validation error for invalid role, got nil")
		}
		if _, ok := errs["role"]; !ok {
			t.Errorf("expected role error, got: %v", errs)
		}
	})
}
