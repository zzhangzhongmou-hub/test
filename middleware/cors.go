package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORS 允许前端跨域访问
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许的来源（开发环境用 *，生产环境指定具体域名）
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// 允许的 HTTP 方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// 允许的请求头
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
		// 允许携带 Cookie（如果用了 JWT 在 Cookie 里）
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		// 预检请求缓存时间
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		// 处理 OPTIONS 预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // No Content
			return
		}

		c.Next()
	}
}
