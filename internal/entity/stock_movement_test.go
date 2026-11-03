package entity

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"web-coffee/pkg/validator"
)

func TestStockMovementJSONMarshaling(t *testing.T) {
	now := time.Now()
	movement := StockMovement{
		ID:          "sm-1",
		ProductID:   "prod-1",
		UserID:      "user-1",
		Type:        MovementTypeAdjustment,
		Quantity:    5,
		StockBefore: 45,
		StockAfter:  40,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	bytes, err := json.Marshal(movement)
	if err != nil {
		t.Fatalf("failed to marshal StockMovement: %v", err)
	}

	jsonStr := string(bytes)
	if strings.Contains(jsonStr, `"notes":`) {
		t.Errorf("expected notes to be omitted when empty, got: %s", jsonStr)
	}
	if strings.Contains(jsonStr, `"product":`) {
		t.Errorf("expected product to be omitted when nil, got: %s", jsonStr)
	}
	if strings.Contains(jsonStr, `"user":`) {
		t.Errorf("expected user to be omitted when nil, got: %s", jsonStr)
	}
}

func TestStockAdjustmentRequestValidation(t *testing.T) {
	v := validator.New()

	t.Run("StockAdjustmentRequest valid in", func(t *testing.T) {
		req := StockAdjustmentRequest{
			Quantity: 20,
			Type:     MovementTypeIn,
			Notes:    "Restock biji kopi dari roastery",
		}
		if errs := v.Validate(req); errs != nil {
			t.Errorf("expected valid, got errors: %v", errs)
		}
	})

	t.Run("StockAdjustmentRequest valid adjustment", func(t *testing.T) {
		req := StockAdjustmentRequest{
			Quantity: 2,
			Type:     MovementTypeAdjustment,
			Notes:    "Susu tumpah saat operasional",
		}
		if errs := v.Validate(req); errs != nil {
			t.Errorf("expected valid, got errors: %v", errs)
		}
	})

	t.Run("StockAdjustmentRequest invalid quantity and type", func(t *testing.T) {
		req := StockAdjustmentRequest{
			Quantity: 0,              // min=1
			Type:     "invalid_type", // oneof=in out adjustment
		}
		errs := v.Validate(req)
		if errs == nil {
			t.Fatal("expected validation errors, got nil")
		}
		if _, ok := errs["quantity"]; !ok {
			t.Errorf("expected error on quantity, got: %v", errs)
		}
		if _, ok := errs["type"]; !ok {
			t.Errorf("expected error on type, got: %v", errs)
		}
	})
}
