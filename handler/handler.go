package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Handler() {
	router := gin.Default()

	// TODO: Fix CORS issue
	// Configure CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // Allow specific origin
	//config.AllowAllOrigins = true //Allow all origins (Not recommended for production)
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

	// TODO: Add signin, signup, and signout routes
	// https://www.sohamkamani.com/golang/session-cookie-authentication/
	router.GET("/signin", signin)
	// router.GET("/test", test)
	router.POST("/signup", signUp)

	// TODO: Add oauth with google
	// https://medium.com/@_RR/google-oauth-2-0-and-golang-4cc299f8c1ed
	router.GET("/auth/google/login", oauthGoogleLogin)
	router.GET("/auth/google/callback", oauthGoogleCallback)

	router.GET("/teams", getTeams)
	// router.GET("/team/:id", getTeam)

	router.POST("/team", postTeam)

	router.Run("localhost:8080")
}
