package handler

import (
	"fmt"

	"github.com/Sigdriv/paskelabyrint-api/db"
	"github.com/Sigdriv/paskelabyrint-api/model"
	"github.com/gin-gonic/gin"
)

func getUser(c *gin.Context) (model.User, error) {
	session, exists := c.Get("session")
	if !exists {
		return model.User{}, fmt.Errorf("Failed to get session")
	}

	conn, err := db.DBConnect()
	if err != nil {
		return model.User{}, fmt.Errorf("Failed to connect to database >> %s", err)
	}
	defer conn.Close()

	userSession, ok := session.(*model.Session)
	if !ok {
		return model.User{}, fmt.Errorf("Failed to cast session >> %s", err)
	}

	var user model.User
	err = conn.QueryRow(c, `Select id, name, email from users where email = $1`, userSession.Email).Scan(
		&user.ID,
		&user.Name,
		&user.Email)

	return user, err
}
