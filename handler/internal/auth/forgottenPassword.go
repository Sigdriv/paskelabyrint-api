package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/Sigdriv/paskelabyrint-api/db"
	"github.com/Sigdriv/paskelabyrint-api/model"
	"github.com/Sigdriv/paskelabyrint-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
)

func ForgottPassword(c *gin.Context) {
	var email model.ForgottPassword

	err := c.ShouldBindJSON(&email)
	if err != nil {
		log.Error("Failed to bind JSON >> ", err)
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request"})
		return
	}

	conn, err := db.DBConnect()
	if err != nil {
		log.Error("Error connecting to database >> ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	defer conn.Close()

	token, err := generateUniqueToken(c, conn)
	if err != nil {
		log.Warn("Error creating unique token >> ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	var user model.User
	err = conn.QueryRow(c, "SELECT id FROM users WHERE email = $1", email.Email).Scan(&user.ID)
	if err != nil {
		log.Warn("User not found or DB error >> ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error quering user from database"})
		return
	}

	var alreadyRequested string
	err = conn.QueryRow(c, "SELECT id FROM resetpassword WHERE created_by_id = $1", user.ID).Scan(&alreadyRequested)

	if alreadyRequested != "" {
		_, err = conn.Exec(c, "DELETE FROM resetpassword WHERE id = $1", alreadyRequested)
		if err != nil {
			log.Warn("Error deleting resetpassword from database >> ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting from database"})
			return
		}
	}

	currentTime := time.Now()
	expirationDate := currentTime.Add(15 * time.Minute)

	_, err = conn.Exec(c, "INSERT INTO resetpassword (token, expires_at, created_at, created_by_id) VALUES ($1, $2, NOW(), $3)",
		token, expirationDate, user.ID,
	)
	if err != nil {
		log.Warn("Error inserting resetpassword into database >> ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting database"})
		return
	}

	domain := os.Getenv("DOMAIN")
	if domain == "" {
		log.Error("Missing DOMAIN environment variable")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal configuration error"})
		return
	}

	body := utils.BuildResetEmailBody(domain, token)

	err = utils.SendEmail(email.Email, "glemt-passord@paskelabyrint.no", "Tilbakestill passord", body)
	if err != nil {
		log.Warn("Error sending email >> ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending email"})
		return
	}

	log.Info("New password requestet successfully sent")
	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
	return
}

func generateUniqueToken(c *gin.Context, conn *pgxpool.Pool) (token string, err error) {
	token = uuid.New().String()

	return tokenCheck(token, c, conn)
}

func tokenCheck(token string, c *gin.Context, conn *pgxpool.Pool) (string, error) {
	var existingToken string
	err := conn.QueryRow(c, "SELECT token FROM resetpassword WHERE token = $1", token).Scan(&existingToken)

	if err != nil || existingToken == "" {
		return token, nil
	}

	return generateUniqueToken(c, conn)
}
