package handler

import "github.com/gin-gonic/gin"

func Handler() {
	router := gin.Default()

	// TODO: Add signin, signup, and signout routes
	// https://www.sohamkamani.com/golang/session-cookie-authentication/

	// TODO: Add oauth with google
	// https://medium.com/@_RR/google-oauth-2-0-and-golang-4cc299f8c1ed

	router.GET("/teams", getTeams)
	// router.GET("/team/:id", getTeam)

	router.POST("/team", postTeam)

	router.Run("localhost:8080")
}
