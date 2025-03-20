package handler

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/Sigdriv/paskelabyrint-api/db"
	"github.com/Sigdriv/paskelabyrint-api/mapping"
	"github.com/Sigdriv/paskelabyrint-api/model"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var sessions = map[string]model.Session{}

func signUp(c *gin.Context) {
	var newUser model.User

	err := c.ShouldBindJSON(&newUser)
	if err != nil {
		log.Error("Failed to bind JSON >> ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	conn, err := db.DBConnect()
	if err != nil {
		log.Error("Error connecting to database >> ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	defer conn.Close()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Error hashing password >> ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	newUser.Password = string(hashedPassword)

	var existingUser string
	err = conn.QueryRow(c, "SELECT email FROM user WHERE email = $1", newUser.Email).Scan(&existingUser)
	if err != nil || existingUser != "" {
		log.Warn("User already exists >> ", newUser.Email)
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	_, err = conn.Exec(c, "INSERT INTO users (name, email, password) VALUES ($1, $2, $3)",
		newUser.Name,
		newUser.Email,
		newUser.Password,
	)
	if err != nil {
		log.Error("Error inserting user into database >> ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	log.Infof("New user added: %+v", newUser)
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func signin(c *gin.Context) {
	var creds model.Creadentials

	err := c.ShouldBindJSON(&creds)
	if err != nil {
		log.Error("Failed to bind JSON >> ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
		return
	}

	conn, err := db.DBConnect()
	if err != nil {
		log.Error("Error connecting to database >> ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	defer conn.Close()

	var expectedPassword string
	err = conn.QueryRow(c, "SELECT password FROM users WHERE email = $1", creds.Email).Scan(&expectedPassword)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Warn("User not found >> ", creds.Email)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		log.Error("Error querying database >> ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(expectedPassword), []byte(creds.Password))

	if err != nil {
		log.Warn("Password mismatch for user >> ", creds.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	sessionToken := generateSecureToken(128)

	expireTime := 24 // Default to 24 hours
	if creds.Remember {
		expireTime = 720 // 30 days
	}
	expiresAt := time.Now().Add(time.Duration(expireTime) * time.Hour)

	_, err = conn.Exec(c, "DELETE FROM Sessions WHERE email = $1", creds.Email)
	if err != nil {
		log.Error("Error deleting old session >> ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	session := model.Session{
		ID:        sessionToken,
		Email:     creds.Email,
		CreatedAt: time.Now(),
		Expiry:    expiresAt,
	}

	dbSession := mapping.SessionToSBSession(session)
	_, err = conn.Exec(c, "INSERT INTO Sessions (id, email, created_at, expires_at) VALUES ($1, $2, $3, $4)",
		dbSession.ID,
		dbSession.Email,
		dbSession.CreatedAt,
		dbSession.Expiry)

	if err != nil {
		log.Error("Error inserting session into database >> ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   true,
	})
}

func generateSecureToken(length int) string {
	token := make([]byte, length)
	_, err := rand.Read(token)
	if err != nil {
		log.Error("Error generating secure token >> ", err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(token)
}

func test(c *gin.Context) {
	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			log.Warn("No session token found >> ", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No session token found"})
			return
		}

		log.Error("Error retrieving session token >> ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	userSession, exisits := sessions[sessionToken]
	if !exisits {
		log.Warn("Session token not found >> ", sessionToken)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session token"})
		return
	}

	if userSession.IsExpired() {
		delete(sessions, sessionToken)
		log.Warn("Session token expired >> ", sessionToken)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session token expired"})
		return
	}

	log.Infof("Session token valid >> %s, User: %s", sessionToken, userSession.Email)
	c.JSON(http.StatusOK, gin.H{"message": "Session token valid", "user": userSession.Email})
}
