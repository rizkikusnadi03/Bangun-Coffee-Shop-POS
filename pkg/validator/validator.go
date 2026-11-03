package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator adalah wrapper untuk *validator.Validate dari go-playground/validator.
type Validator struct {
	validate *validator.Validate
}

// New mengembalikan instance Validator baru yang telah dikonfigurasi.
func New() *Validator {
	v := validator.New()

	// Mengatur tag name function bawaan agar membaca JSON tag jika tersedia
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validator{
		validate: v,
	}
}

// Validate menjalankan validasi struct input dan mengembalikan map error dalam format snake_case JSON tag dan pesan bahasa Indonesia.
func (v *Validator) Validate(i interface{}) map[string]string {
	err := v.validate.Struct(i)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return map[string]string{
			"error": err.Error(),
		}
	}

	errMap := make(map[string]string, len(validationErrors))
	for _, e := range validationErrors {
		fieldName := getJSONFieldName(i, e)
		errMap[fieldName] = formatErrorMessage(e)
	}

	return errMap
}

// getJSONFieldName mengekstrak nama field dalam format snake_case dari JSON tag struct melalui refleksi.
func getJSONFieldName(i interface{}, e validator.FieldError) string {
	structField := e.StructField()
	if structField == "" {
		structField = e.Field()
	}

	if i != nil && structField != "" {
		t := reflect.TypeOf(i)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() == reflect.Struct {
			if f, ok := t.FieldByName(structField); ok {
				jsonTag := f.Tag.Get("json")
				if jsonTag != "" && jsonTag != "-" {
					parts := strings.SplitN(jsonTag, ",", 2)
					if parts[0] != "" {
						return parts[0]
					}
				}
			}
		}
	}

	if e.Field() != "" {
		return e.Field()
	}

	return structField
}

// formatErrorMessage memformat pesan kesalahan validasi ke dalam Bahasa Indonesia.
func formatErrorMessage(e validator.FieldError) string {
	param := e.Param()
	kind := e.Kind()

	switch e.Tag() {
	case "required":
		return "wajib diisi"
	case "min":
		if isNumberKind(kind) {
			return fmt.Sprintf("minimal %s", param)
		}
		return fmt.Sprintf("minimal %s karakter", param)
	case "max":
		if isNumberKind(kind) {
			return fmt.Sprintf("maksimal %s", param)
		}
		return fmt.Sprintf("maksimal %s karakter", param)
	case "email":
		return "format email tidak valid"
	case "oneof":
		return fmt.Sprintf("harus salah satu dari: %s", param)
	case "uuid4":
		return "format UUID tidak valid"
	case "len":
		return fmt.Sprintf("harus tepat %s karakter", param)
	case "numeric":
		return "harus berupa angka"
	default:
		return "tidak valid"
	}
}

// isNumberKind mengecek apakah tipe data field merupakan numerik (integer atau float).
func isNumberKind(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}
