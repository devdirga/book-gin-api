package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   200,
			"message":  "pong",
			"is_error": false,
			"data": gin.H{
				"reference_number": "{{reference_number}}",
				"channel_id":       "{{channel_id}}",
				"info":             "randomstring",
				"otp_type":         0,
				"token":            "{{register_token}}",
				"otp_channel":      0,
			},
		})
	})
	r.Run()
}
