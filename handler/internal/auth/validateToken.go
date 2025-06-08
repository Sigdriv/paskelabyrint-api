package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ValidateToken(c *gin.Context) {
	token := c.Param("token")

	valid, err := validateToken(token, c)
	if err != nil {
		logrus.Error("Token validation failed >> ", err)
		c.JSON(http.StatusGone, gin.H{"error": err.Error()})
		return
	}
	if !valid {
		logrus.Warn("Token is not valid or expired >> ", token)
		c.JSON(http.StatusGone, gin.H{"error": "Token is not valid or expired"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true})
	logrus.Info("Token validation successful for token: ", token)
}
