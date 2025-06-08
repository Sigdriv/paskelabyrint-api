package handler

import (
	"log"
	"net"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct{}

func CreateHandler() (srv Handler, err error) {

	return
}

func (srv *Handler) CreateGinGroup() {
	router := gin.Default()

	router.Use(configureCors())

	// Public routes
	public := router.Group("/")

	public.POST("/auth/login", signin)
	public.POST("/auth/signup", signUp)
	public.POST("/auth/forgott-password", srv.handleForgottenPassword)
	public.GET("/auth/google/login", oauthGoogleLogin)
	public.GET("/auth/google/callback", oauthGoogleCallback)
	public.GET("/auth/validate-token/:token", srv.handleValidateToken)
	public.POST("/auth/reset-password/:token", srv.handleResetPassword)

	// Protected routes
	protected := router.Group("/")
	protected.Use(AuthMiddleware())

	protected.GET("/user", userBySession)
	protected.GET("/auth/logout", logOut)

	protected.GET("/teams", getTeams)
	protected.POST("/team", postTeam)

	router.Run("localhost:8080")
}

func configureCors() gin.HandlerFunc {
	// local := fmt.Sprintf("http://%s:3000", getLocalIP())
	local := "http://127.0.0.1:3000"

	return cors.New(cors.Config{
		AllowOrigins:     []string{local, "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

func getLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP
}
