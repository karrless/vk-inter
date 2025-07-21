package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseQueryInt(c *gin.Context, key string, defaultValue int) int {
	if v := c.Query(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultValue
}

func ParseQueryFloat(c *gin.Context, key string, defaultValue float64) float64 {
	if v := c.Query(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return defaultValue
}
