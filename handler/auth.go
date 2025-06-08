package handler

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Sigdriv/paskelabyrint-api/db"
	"github.com/Sigdriv/paskelabyrint-api/handler/internal/auth"

	"github.com/Sigdriv/paskelabyrint-api/model"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func (*Handler) handleForgottenPassword(c *gin.Context) {
	auth.ForgottPassword(c)
}

func (*Handler) handleValidateToken(c *gin.Context) {
	auth.ValidateToken(c)
}

func (*Handler) handleResetPassword(c *gin.Context) {
	auth.ResetPassword(c)
}

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
	err = conn.QueryRow(c, "SELECT email FROM users WHERE email = $1", newUser.Email).Scan(&existingUser)
	if existingUser != "" {
		log.Warn("User already exists >> ", err)
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

	setCookie(c, conn, creds.Email, creds.Remember, true)
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

// func test(c *gin.Context) {
// 	sessionToken, err := c.Cookie("session_token")
// 	if err != nil {
// 		if err == http.ErrNoCookie {
// 			log.Warn("No session token found >> ", err)
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "No session token found"})
// 			return
// 		}

// 		log.Error("Error retrieving session token >> ", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
// 		return
// 	}

// 	userSession, exisits := sessions[sessionToken]
// 	if !exisits {
// 		log.Warn("Session token not found >> ", sessionToken)
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session token"})
// 		return
// 	}

// 	if userSession.IsExpired() {
// 		delete(sessions, sessionToken)
// 		log.Warn("Session token expired >> ", sessionToken)
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session token expired"})
// 		return
// 	}

// 	log.Infof("Session token valid >> %s, User: %s", sessionToken, userSession.Email)
// 	c.JSON(http.StatusOK, gin.H{"message": "Session token valid", "user": userSession.Email})
// }

func generateSecureTokenWithExpiry(length int, expiry time.Time, userRole model.Role) string {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Error("Error generating secure token >> ", err)
		return ""
	}

	randomPart := base64.StdEncoding.EncodeToString(randomBytes)

	expiryUnix := expiry.Unix()

	combinedToken := fmt.Sprintf("%s|%d|%s", randomPart, expiryUnix, userRole)

	return base64.StdEncoding.EncodeToString([]byte(combinedToken))
}

func setCookie(c *gin.Context, conn *pgxpool.Pool, email string, remember bool, credentials bool) {
	expireTime := 24 // Default to 24 hours
	if remember {
		expireTime = 720 // 30 days
	}
	expiresAt := time.Now().Add(time.Duration(expireTime) * time.Hour)

	var userRole model.Role
	err := conn.QueryRow(c, "SELECT role FROM users WHERE email = $1", email).Scan(&userRole)

	sessionToken := generateSecureTokenWithExpiry(32, expiresAt, userRole)

	_, err = conn.Exec(c, "DELETE FROM Sessions WHERE email = $1", email)
	if err != nil {
		log.Error("Error deleting old session >> ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	token, _, err := extractTokenData(sessionToken)
	if err != nil {
		log.Error("Error extracting token data >> ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	session := model.Session{
		ID:        token,
		Email:     email,
		CreatedAt: time.Now(),
		Expiry:    expiresAt,
	}

	_, err = conn.Exec(c, "INSERT INTO Sessions (id, email, created_at, expires_at) VALUES ($1, $2, $3, $4)",
		session.ID,
		session.Email,
		session.CreatedAt,
		session.Expiry)

	if err != nil {
		log.Error("Error inserting session into database >> ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.SetCookie(
		"session_token",
		sessionToken,
		int(expireTime*60*60),
		"/",
		"",
		false,
		true,
	)

	redirectUrl := os.Getenv("REDIRECT_URL")
	if redirectUrl == "" {
		redirectUrl = "https://paskelabyrint.no/dev"
	}

	log.Infof("Session token set >> %s, User: %s", sessionToken, email)
	if credentials {
		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "session_token": sessionToken})
	} else {
		c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
	}
}

func extractTokenData(token string) (string, time.Time, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		log.Error("Error decoding token >> ", err)
		return "", time.Time{}, err
	}

	parts := strings.Split(string(decodedBytes), "|")
	if len(parts) != 3 {
		log.Error("Invalid token format >> ", token)
		return "", time.Time{}, fmt.Errorf("invalid token format")
	}

	expiryUnix, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		log.Error("Error parsing expiry time >> ", err)
		return "", time.Time{}, err
	}

	expiryTime := time.Unix(expiryUnix, 0)

	return parts[0], expiryTime, nil
}

func validateSessionToken(c *gin.Context, conn *pgxpool.Pool) (*model.Session, error) {
	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			log.Warn("No session token found >> ", err)
			return nil, fmt.Errorf("no session token found")
		}
		log.Error("Error retrieving session token >> ", err)
		return nil, fmt.Errorf("error retrieving session token")
	}

	token, expiry, err := extractTokenData(sessionToken)
	if err != nil {
		log.Error("Error extracting token data >> ", err)
		return nil, fmt.Errorf("error extracting token data")
	}

	if time.Now().After(expiry) {
		log.Warn("Session token expired >> ", sessionToken)
		return nil, fmt.Errorf("session token expired")
	}

	var session model.Session
	err = conn.QueryRow(c, "SELECT id, email, created_at, expires_at FROM Sessions WHERE id = $1", token).Scan(
		&session.ID,
		&session.Email,
		&session.CreatedAt,
		&session.Expiry,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Warn("Token not found in database >> ", token)
			return nil, fmt.Errorf("Token not found in database >> %s", err)
		}
		log.Error("Error querying database >> ", err)
		return nil, fmt.Errorf("error querying database")
	}

	if session.Expiry.Before(time.Now()) {
		log.Warn("Session token expired >> ", sessionToken)
		return nil, fmt.Errorf("session token expired")
	}
	log.Infof("Session token valid >> %s, User: %s", sessionToken, session.Email)
	return &session, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := db.DBConnect()
		defer conn.Close()
		if err != nil {
			log.Error("Error connecting to database >> ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		session, err := validateSessionToken(c, conn)
		if err != nil {
			log.Warn("Session validation failed >> ", err)
			c.SetCookie(
				"session_token",
				"Cookie cler",
				-1,
				"/",
				"",
				false,
				true,
			)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Set("session", session)
		log.Infof("Session validated >> %s, User: %s", session.ID, session.Email)
		c.Next()
		// Optionally, you can add logic here to refresh the session token if needed
		// For example, you can check if the session is about to expire and generate a new token
	}
}

func userBySession(c *gin.Context) {
	session, exists := c.Get("session")
	if !exists {
		log.Warn("Session not found in context >> ")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userSession, ok := session.(*model.Session)
	if !ok {
		log.Error("Invalid session type >> ")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	conn, err := db.DBConnect()
	if err != nil {
		log.Error("Error connecting to database >> ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	defer conn.Close()

	var user model.User
	err = conn.QueryRow(c, "select id, name, email, created_at, role, avatar from users where email = $1", userSession.Email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.Role,
		&user.Avatar,
	)

	log.Infof("User session retrieved >> %s", user)
	c.JSON(http.StatusOK, user)
}

func logOut(c *gin.Context) {
	session, exists := c.Get("session")
	if !exists {
		log.Warn("Session not found in context >> ")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userSession, ok := session.(*model.Session)
	if !ok {
		log.Error("Invalid session type >> ")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	conn, err := db.DBConnect()
	if err != nil {
		log.Error("Error connecting to database >> ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	defer conn.Close()

	_, err = conn.Exec(c, "DELETE FROM sessions WHERE id = $1", userSession.ID)
	if err != nil {
		log.Error("Error deleting session from database >> ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	redirectUrl := os.Getenv("REDIRECT_URL")
	if redirectUrl == "" {
		redirectUrl = "https://paskelabyrint.no"
	}

	c.SetCookie("session_token", "", -1, "/", "", false, true)
	// c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
	c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
}
