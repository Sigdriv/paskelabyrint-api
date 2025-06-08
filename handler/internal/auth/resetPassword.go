package auth

import (
	"net/http"

	"github.com/Sigdriv/paskelabyrint-api/db"
	"github.com/Sigdriv/paskelabyrint-api/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ResetPassword(c *gin.Context) {

	token := c.Param("token")
	valid, err := validateToken(token, c)
	if err != nil {
		c.JSON(http.StatusGone, gin.H{"error": "Token validation failed"})
		return
	}
	if !valid {
		c.JSON(http.StatusGone, gin.H{"error": "Token is not valid or expired"})
		return
	}

	conn, err := db.DBConnect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}
	defer conn.Close()

	var resetPassword model.ResetPassword
	err = c.ShouldBindJSON(&resetPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if resetPassword.Password != resetPassword.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(resetPassword.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	_, err = conn.Exec(c, "UPDATE users SET password = $1 WHERE id = (SELECT created_by_id FROM resetpassword WHERE token = $2)", hashedPassword, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating password in database"})
		return
	}
	_, err = conn.Exec(c, "DELETE FROM resetpassword WHERE token = $1", token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting reset password token from database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}
