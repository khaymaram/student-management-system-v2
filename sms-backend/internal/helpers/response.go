package helpers

import "github.com/gin-gonic/gin"

// SuccessResponse wraps a successful API result in the standard { success, data } shape.
func SuccessResponse(
	c *gin.Context,
	status int,
	data interface{},
) {

	c.JSON(status, gin.H{
		"success": true,
		"data":    data,
	})
}

// ErrorResponse wraps a failed API result in the standard { success, error } shape.
func ErrorResponse(
	c *gin.Context,
	status int,
	message string,
) {

	c.JSON(status, gin.H{
		"success": false,
		"error":   message,
	})
}
