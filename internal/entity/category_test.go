package entity

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"web-coffee/pkg/validator"
)

func TestCategoryJSONMarshaling(t *testing.T) {
	now := time.Now()
	cat := Category{
		ID:        "cat-1",
		Name:      "Espresso Based",
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: &now,
	}

	bytes, err := json.Marshal(cat)
	if err != nil {
		t.Fatalf("failed to marshal Category: %v", err)
	}

	jsonStr := string(bytes)
	if strings.Contains(jsonStr, "deleted_at") {
		t.Errorf("expected deleted_at to be omitted from JSON, got: %s", jsonStr)
	}
	if !strings.Contains(jsonStr, `"name":"Espresso Based"`) {
		t.Errorf("expected name to be in JSON, got: %s", jsonStr)
	}
}

func TestCategoryRequestsValidation(t *testing.T) {
	v := validator.New()

	t.Run("CreateCategoryRequest valid", func(t *testing.T) {
		req := CreateCategoryRequest{Name: "Manual Brew"}
		if errs := v.Validate(req); errs != nil {
			t.Errorf("expected valid, got: %v", errs)
		}
	})

	t.Run("CreateCategoryRequest invalid", func(t *testing.T) {
		req := CreateCategoryRequest{Name: "A"} // too short (< 2)
		errs := v.Validate(req)
		if errs == nil {
			t.Fatal("expected validation error, got nil")
		}
		if _, ok := errs["name"]; !ok {
			t.Errorf("expected error on name, got: %v", errs)
		}
	})

	t.Run("UpdateCategoryRequest valid", func(t *testing.T) {
		req := UpdateCategoryRequest{Name: "Tea & Non-Coffee"}
		if errs := v.Validate(req); errs != nil {
			t.Errorf("expected valid, got: %v", errs)
		}
	})

	t.Run("UpdateCategoryRequest empty", func(t *testing.T) {
		req := UpdateCategoryRequest{Name: ""}
		errs := v.Validate(req)
		if errs == nil {
			t.Fatal("expected validation error on empty name, got nil")
		}
		if _, ok := errs["name"]; !ok {
			t.Errorf("expected error on name, got: %v", errs)
		}
	})
}
