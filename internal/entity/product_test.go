package entity

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"web-coffee/pkg/validator"
)

func TestProductPriceInRupiah(t *testing.T) {
	prod := Product{
		ID:    "prod-1",
		Name:  "Iced Caffe Latte",
		Price: 2800000, // Rp 28.000,00
	}

	expectedRupiah := 28000.0
	if prod.PriceInRupiah() != expectedRupiah {
		t.Errorf("expected %.2f Rupiah, got %.2f", expectedRupiah, prod.PriceInRupiah())
	}
}

func TestProductJSONMarshaling(t *testing.T) {
	now := time.Now()
	prod := Product{
		ID:         "prod-1",
		CategoryID: "cat-1",
		Name:       "Espresso",
		Price:      1800000,
		Stock:      50,
		IsActive:   true,
		CreatedAt:  now,
		UpdatedAt:  now,
		DeletedAt:  &now,
	}

	bytes, err := json.Marshal(prod)
	if err != nil {
		t.Fatalf("failed to marshal product: %v", err)
	}

	jsonStr := string(bytes)
	if strings.Contains(jsonStr, "deleted_at") {
		t.Errorf("expected deleted_at to be omitted from JSON, got: %s", jsonStr)
	}
	if strings.Contains(jsonStr, `"category":`) {
		t.Errorf("expected category to be omitted when nil, got: %s", jsonStr)
	}
}

func TestProductRequestsValidation(t *testing.T) {
	v := validator.New()

	t.Run("CreateProductRequest valid", func(t *testing.T) {
		req := CreateProductRequest{
			CategoryID:  "e1cfbfb0-775c-4c32-a8a9-8143c33c2952",
			Name:        "Iced Caramel Macchiato",
			Description: "Espresso dengan vanila dan saus karamel",
			Price:       3200000,
			ImageURL:    "https://example.com/macchiato.png",
		}
		if errs := v.Validate(req); errs != nil {
			t.Errorf("expected valid, got: %v", errs)
		}
	})

	t.Run("CreateProductRequest invalid", func(t *testing.T) {
		req := CreateProductRequest{
			CategoryID: "not-a-uuid",
			Name:       "A", // too short (< 2)
			Price:      50,  // below 100 sen
			ImageURL:   "invalid-url",
		}
		errs := v.Validate(req)
		if errs == nil {
			t.Fatal("expected validation errors, got nil")
		}
		if _, ok := errs["category_id"]; !ok {
			t.Errorf("expected error on category_id, got: %v", errs)
		}
		if _, ok := errs["name"]; !ok {
			t.Errorf("expected error on name, got: %v", errs)
		}
		if _, ok := errs["price"]; !ok {
			t.Errorf("expected error on price, got: %v", errs)
		}
		if _, ok := errs["image_url"]; !ok {
			t.Errorf("expected error on image_url, got: %v", errs)
		}
	})

	t.Run("UpdateProductRequest valid partial", func(t *testing.T) {
		isActive := false
		req := UpdateProductRequest{
			Name:     "Updated Name",
			IsActive: &isActive,
		}
		if errs := v.Validate(req); errs != nil {
			t.Errorf("expected valid for partial update, got: %v", errs)
		}
	})
}
