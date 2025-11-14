package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"kami/internal/config"
	"kami/internal/logger"
	"kami/internal/utils"
)

func EnsureIdentity() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		email := ""
		valid := false
		if err == nil && token != "" {
			if em, exp, perr := utils.ParseToken(config.Config.JWT_SECRET, token); perr == nil && time.Now().Before(exp) {
				email = em
				valid = true
			}
		}
		if !valid {
			generated := false
			if em, ok := utils.GenerateRandomEmail(); ok {
				email = em
				generated = true
			} else {
				email = config.Config.STMP_USER
			}
			tok, exp, _ := utils.GenerateToken(config.Config.JWT_SECRET, email, 24*time.Hour)
			maxAge := int(time.Until(exp).Seconds())
			c.SetCookie("token", tok, maxAge, "/", "", false, true)
			logger.LogRequest(c, "assign_email", email, map[string]any{
				"source":    "middleware",
				"generated": generated,
			})
		}
		c.Set("email", email)
		c.Next()
	}
}
