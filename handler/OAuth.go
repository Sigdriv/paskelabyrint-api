package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Sigdriv/paskelabyrint-api/db"
	"github.com/Sigdriv/paskelabyrint-api/model"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func oauthGoogleLogin(c *gin.Context) {
	oauthState := generateSecureToken(200)

	var googleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("AUTH_GOOGLE_CALLBACK_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	url := googleOauthConfig.AuthCodeURL(oauthState, oauth2.AccessTypeOffline)

	c.Redirect(http.StatusTemporaryRedirect, url)
}

func oauthGoogleCallback(c *gin.Context) {
	data, err := getUserDataFromGoogle(c.Request.FormValue("state"), c.Request.FormValue("code"))
	if err != nil {
		log.Error("Failed to get user data from Google >> ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user data"})
		return
	}

	var userInfo model.GoogleUserInfo
	err = json.Unmarshal(data, &userInfo)
	if err != nil {
		log.Error("Failed to unmarshal user data >> ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal user data"})
		return
	}

	conn, err := db.DBConnect()
	if err != nil {
		log.Error("Database connection error >> ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}

	defer conn.Close()

	var userExists string
	conn.QueryRow(context.Background(), "SELECT id FROM users WHERE email = $1", userInfo.Email).Scan(&userExists)
	if userExists != "" {
		log.Info("User already exists >> ", userInfo.Email)
	}

	if userExists == "" {
		err = saveUser(conn, userInfo)
		if err != nil {
			log.Error("Failed to save user data >> ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user data"})

			return
		}
	}

	if userExists != "" {
		err = updateUser(conn, userInfo, userExists)
		if err != nil {
			log.Error("Failed to update user data >> ", err)
		}
	}

	setCookie(c, conn, userInfo.Email, false, false)
}

func getUserDataFromGoogle(state string, code string) ([]byte, error) {
	// TODO: Validate state with oauthState
	if state != state {
		return nil, fmt.Errorf("Invalid OAuth state")
	}

	var googleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("AUTH_GOOGLE_CALLBACK_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("Failed to exchange token >> %v", err)
	}

	const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to get user data >> %v", err)
	}
	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body >> %v", err)
	}

	return contents, nil
}

func saveUser(conn *pgxpool.Pool, userInfo model.GoogleUserInfo) error {

	_, err := conn.Exec(context.Background(),
		"INSERT INTO users(name, email, avatar) VALUES ($1, $2, $3)",
		userInfo.Name, userInfo.Email, userInfo.Picture)
	if err != nil {
		return fmt.Errorf("Failed to insert user into database >> %v", err)
	}

	return nil
}

func updateUser(conn *pgxpool.Pool, userInfo model.GoogleUserInfo, userIndex string) error {

	_, err := conn.Exec(context.Background(),
		"UPDATE users SET name = $1, avatar = $2 WHERE id = $3",
		userInfo.Name, userInfo.Picture, userIndex)
	if err != nil {
		return fmt.Errorf("Failed to update user into database >> %v", err)
	}

	return nil
}
