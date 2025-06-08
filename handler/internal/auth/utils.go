package auth

import (
	"fmt"
	"time"

	"github.com/Sigdriv/paskelabyrint-api/db"
	"github.com/Sigdriv/paskelabyrint-api/model"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

func validateToken(token string, c *gin.Context) (valid bool, err error) {

	if token == "" {
		err = fmt.Errorf("Missing token in request")
		return
	}

	conn, err := db.DBConnect()
	if err != nil {
		err = fmt.Errorf("Error connecting to database >> %v ", err)
		return
	}
	defer conn.Close()

	var resetPassword model.ResetPassword
	err = conn.QueryRow(c, "SELECT expires_at FROM resetpassword WHERE token = $1", token).Scan(&resetPassword.ExpiresAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = fmt.Errorf("Token not found in database")
			return
		}
		err = fmt.Errorf("Error querying reset password token >> %v", err)
		return
	}

	if resetPassword.ExpiresAt.Before(time.Now()) {
		err = fmt.Errorf("Token is expired")
		return
	}

	valid = true
	logrus.Info("Token validation successful for token: ", token)
	return
}
