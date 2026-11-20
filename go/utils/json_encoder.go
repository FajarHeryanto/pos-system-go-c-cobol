// ==================================================================
// utils/json_encoder.go — JSON Response Helper
// Sistem Informasi Penjualan (Point of Sales)
// ==================================================================

package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"pos-system/models"
)

// RespondSuccess mengirim JSON response sukses dengan data.
func RespondSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, models.APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// RespondError mengirim JSON response error.
func RespondError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, models.APIResponse{
		Status:  "error",
		Message: message,
	})
}

// RespondValidationError mengirim JSON response error validasi dengan detail per-field.
func RespondValidationError(c *gin.Context, errors []models.FieldError) {
	c.JSON(http.StatusBadRequest, models.APIResponse{
		Status:  "error",
		Message: "Validasi gagal",
		Errors:  errors,
	})
}
