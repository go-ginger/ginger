package ginger

import "github.com/gin-gonic/gin"

func CORS(c *gin.Context) {
	if config.CorsEnabled {
		c.Writer.Header().Set("Access-Control-Allow-Origin", config.CorsAllowOrigins)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", config.CorsAllowCredentials)
		c.Writer.Header().Set("Access-Control-Allow-Headers", config.CorsAllowHeaders)
		c.Writer.Header().Set("Access-Control-Allow-Methods", config.CorsAllowMethods)
	}
	c.Next()
}
