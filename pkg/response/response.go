package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response merepresentasikan format standar respon API JSON.
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// Meta merepresentasikan informasi paginasi pada respon data list.
type Meta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// Success mengirimkan respon JSON sukses kustom dengan status code, pesan, dan data.
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created mengirimkan respon sukses 201 Created dengan data entitas yang baru dibuat.
func Created(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusCreated, message, data)
}

// OK mengirimkan respon sukses 200 OK dengan data payload.
func OK(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusOK, message, data)
}

// Paginated mengirimkan respon sukses 200 OK untuk daftar data yang disertai metadata paginasi.
func Paginated(c *gin.Context, message string, data interface{}, meta Meta) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    &meta,
	})
}

// Error mengirimkan respon error umum dengan status code dan pesan.
func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
	})
}

// BadRequest mengirimkan respon error 400 Bad Request.
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

// Unauthorized mengirimkan respon error 401 Unauthorized (token tidak ada atau tidak valid).
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message)
}

// Forbidden mengirimkan respon error 403 Forbidden (akses ditolak untuk role saat ini).
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message)
}

// NotFound mengirimkan respon error 404 Not Found (data atau resource tidak ditemukan).
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}

// InternalError mengirimkan respon error 500 Internal Server Error.
func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
}

// ValidationError mengirimkan respon error 422 Unprocessable Entity dengan rincian kegagalan validasi field.
func ValidationError(c *gin.Context, errors interface{}) {
	c.JSON(http.StatusUnprocessableEntity, Response{
		Success: false,
		Message: "Validasi gagal",
		Errors:  errors,
	})
}
