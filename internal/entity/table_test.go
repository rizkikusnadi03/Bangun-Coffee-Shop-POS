package entity

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"web-coffee/pkg/validator"
)

func TestTableJSONMarshaling(t *testing.T) {
	now := time.Now()
	tbl := Table{
		ID:        "tbl-1",
		Name:      "Meja 01",
		Capacity:  4,
		Status:    TableStatusAvailable,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: &now,
	}

	bytes, err := json.Marshal(tbl)
	if err != nil {
		t.Fatalf("failed to marshal Table: %v", err)
	}

	jsonStr := string(bytes)
	if strings.Contains(jsonStr, "deleted_at") {
		t.Errorf("expected deleted_at to be omitted from JSON, got: %s", jsonStr)
	}
	if !strings.Contains(jsonStr, `"status":"available"`) {
		t.Errorf("expected status available in JSON, got: %s", jsonStr)
	}
}

func TestTableRequestsValidation(t *testing.T) {
	v := validator.New()

	t.Run("CreateTableRequest valid", func(t *testing.T) {
		req := CreateTableRequest{
			Name:     "Meja 02 (Outdoor)",
			Capacity: 4,
		}
		if errs := v.Validate(req); errs != nil {
			t.Errorf("expected valid, got: %v", errs)
		}
	})

	t.Run("CreateTableRequest invalid", func(t *testing.T) {
		req := CreateTableRequest{
			Name:     "",
			Capacity: 0,
		}
		errs := v.Validate(req)
		if errs == nil {
			t.Fatal("expected validation errors, got nil")
		}
		if _, ok := errs["name"]; !ok {
			t.Errorf("expected error on name, got: %v", errs)
		}
		if _, ok := errs["capacity"]; !ok {
			t.Errorf("expected error on capacity, got: %v", errs)
		}
	})

	t.Run("CreateTableRequest capacity exceeds max", func(t *testing.T) {
		req := CreateTableRequest{
			Name:     "Meja Banquet",
			Capacity: 60, // max=50
		}
		errs := v.Validate(req)
		if errs == nil {
			t.Fatal("expected validation error, got nil")
		}
		if _, ok := errs["capacity"]; !ok {
			t.Errorf("expected error on capacity, got: %v", errs)
		}
	})

	t.Run("UpdateTableRequest valid partial", func(t *testing.T) {
		req := UpdateTableRequest{
			Status: TableStatusOccupied,
		}
		if errs := v.Validate(req); errs != nil {
			t.Errorf("expected valid partial update, got: %v", errs)
		}
	})

	t.Run("UpdateTableRequest invalid status", func(t *testing.T) {
		req := UpdateTableRequest{
			Status: "broken",
		}
		errs := v.Validate(req)
		if errs == nil {
			t.Fatal("expected validation error on status, got nil")
		}
		if _, ok := errs["status"]; !ok {
			t.Errorf("expected error on status, got: %v", errs)
		}
	})
}
