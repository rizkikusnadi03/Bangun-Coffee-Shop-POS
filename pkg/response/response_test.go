package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupContext() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestResponseSuccessHelpers(t *testing.T) {
	t.Run("OK helper", func(t *testing.T) {
		c, w := setupContext()
		OK(c, "success message", map[string]string{"foo": "bar"})

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}

		var resp Response
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if !resp.Success || resp.Message != "success message" {
			t.Errorf("unexpected response: %+v", resp)
		}
		if resp.Meta != nil || resp.Errors != nil {
			t.Errorf("expected Meta and Errors to be omitted, got %+v", resp)
		}
	})

	t.Run("Created helper", func(t *testing.T) {
		c, w := setupContext()
		Created(c, "item created", map[string]int{"id": 123})

		if w.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d", w.Code)
		}

		var resp Response
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if !resp.Success || resp.Message != "item created" {
			t.Errorf("unexpected response: %+v", resp)
		}
	})

	t.Run("Paginated helper", func(t *testing.T) {
		c, w := setupContext()
		meta := Meta{Page: 1, PerPage: 10, Total: 25, TotalPages: 3}
		Paginated(c, "items fetched", []string{"a", "b"}, meta)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}

		var resp Response
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if !resp.Success || resp.Meta == nil || resp.Meta.Total != 25 {
			t.Errorf("unexpected paginated response: %+v", resp)
		}
	})
}

func TestResponseErrorHelpers(t *testing.T) {
	cases := []struct {
		name       string
		call       func(*gin.Context)
		expectCode int
	}{
		{"BadRequest", func(c *gin.Context) { BadRequest(c, "bad request") }, http.StatusBadRequest},
		{"Unauthorized", func(c *gin.Context) { Unauthorized(c, "unauthorized") }, http.StatusUnauthorized},
		{"Forbidden", func(c *gin.Context) { Forbidden(c, "forbidden") }, http.StatusForbidden},
		{"NotFound", func(c *gin.Context) { NotFound(c, "not found") }, http.StatusNotFound},
		{"InternalError", func(c *gin.Context) { InternalError(c, "server error") }, http.StatusInternalServerError},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c, w := setupContext()
			tc.call(c)

			if w.Code != tc.expectCode {
				t.Errorf("expected status %d, got %d", tc.expectCode, w.Code)
			}

			var resp Response
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}
			if resp.Success != false {
				t.Errorf("expected success false, got %v", resp.Success)
			}
			if resp.Data != nil || resp.Meta != nil || resp.Errors != nil {
				t.Errorf("expected Data, Meta, Errors to be omitted, got %+v", resp)
			}
		})
	}

	t.Run("ValidationError helper", func(t *testing.T) {
		c, w := setupContext()
		errs := map[string]string{
			"email": "Format email tidak valid",
		}
		ValidationError(c, errs)

		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected 422, got %d", w.Code)
		}

		var resp Response
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if resp.Success != false || resp.Message != "Validasi gagal" || resp.Errors == nil {
			t.Errorf("unexpected validation response: %+v", resp)
		}
	})
}
