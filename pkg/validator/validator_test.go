package validator

import (
	"strings"
	"testing"
)

type TestInput struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
	Price int64  `json:"price" validate:"required,min=1"`
}

func TestValidateInputValid(t *testing.T) {
	v := New()

	input := TestInput{
		Name:  "Kopi Susu Gula Aren",
		Email: "customer@example.com",
		Price: 18000,
	}

	errs := v.Validate(input)
	if errs != nil {
		t.Errorf("expected nil for valid input, got: %v", errs)
	}
}

func TestValidateAllFieldsEmpty(t *testing.T) {
	v := New()

	input := TestInput{
		Name:  "",
		Email: "",
		Price: 0,
	}

	errs := v.Validate(input)
	if errs == nil {
		t.Fatal("expected validation errors for empty fields, got nil")
	}

	expectedKeys := []string{"name", "email", "price"}
	if len(errs) != len(expectedKeys) {
		t.Errorf("expected map with %d keys, got %d keys: %v", len(expectedKeys), len(errs), errs)
	}

	for _, key := range expectedKeys {
		if _, exists := errs[key]; !exists {
			t.Errorf("expected key %q in error map, got: %v", key, errs)
		}
	}
}

func TestValidateInvalidEmailFormat(t *testing.T) {
	v := New()

	input := TestInput{
		Name:  "Latte",
		Email: "not-an-email",
		Price: 25000,
	}

	errs := v.Validate(input)
	if errs == nil {
		t.Fatal("expected validation error for invalid email, got nil")
	}

	msg, exists := errs["email"]
	if !exists {
		t.Fatalf("expected error for 'email' field, got: %v", errs)
	}

	expected := "format email tidak valid"
	if !strings.Contains(msg, expected) {
		t.Errorf("expected error message to contain %q, got %q", expected, msg)
	}
}

func TestValidateNameTooShort(t *testing.T) {
	v := New()

	input := TestInput{
		Name:  "A", // 1 karakter, syarat min=2
		Email: "latte@example.com",
		Price: 20000,
	}

	errs := v.Validate(input)
	if errs == nil {
		t.Fatal("expected validation error for name too short, got nil")
	}

	msg, exists := errs["name"]
	if !exists {
		t.Fatalf("expected error for 'name' field, got: %v", errs)
	}

	expected := "minimal 2 karakter"
	if !strings.Contains(msg, expected) {
		t.Errorf("expected error message to contain %q, got %q", expected, msg)
	}
}
