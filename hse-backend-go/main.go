package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"hse-backend-go/db"
	"hse-backend-go/handlers"
	"hse-backend-go/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	db.Connect()

	router := gin.Default()
	router.Use(cors.Default()) // for development; restrict origin in production

	router.POST("/api/sign-in", handlers.SignIn)
	router.POST("/api/sign-up", handlers.SignUp)
	router.POST("/api/sign-out", handlers.SignOut)
	router.GET("/api/verify-email", handlers.VerifyEmail)
	router.POST("/api/oauth/microsoft", handlers.MicrosoftOAuthSignIn)

	protected := router.Group("/api")
	protected.Use(middleware.RequireAuth)
	{
		protected.GET("/auth/me", handlers.Me)
		protected.GET("/menu", handlers.GetMyMenu)

		protected.POST("/organizations", handlers.CreateOrganization)
		protected.GET("/organizations/mine", handlers.GetMyOrganization)

		reviewOnly := protected.Group("/organizations")
		reviewOnly.Use(middleware.RequireRole("Admin", "ANP HSE"))
		{
			reviewOnly.GET("", handlers.ListOrganizations)
			reviewOnly.GET("/:id", handlers.GetOrganization)
			reviewOnly.PATCH("/:id/decision", handlers.DecideOrganization)
		}

		adminOnly := protected.Group("")
		adminOnly.Use(middleware.RequireRole("Admin"))
		{
			adminOnly.GET("/users", handlers.ListUsers)
			adminOnly.GET("/roles", handlers.ListRoles)
			adminOnly.PATCH("/users/:id/role", handlers.UpdateUserRole)
			adminOnly.PATCH("/users/:id/status", handlers.UpdateUserStatus)
		}
	}

	api := router.Group("/api/hse")
	{
		api.GET("/submissions", handlers.ListSubmissions)
		api.GET("/submissions/:id", handlers.GetSubmission)
		api.POST("/submissions", handlers.CreateSubmission)
		api.PATCH("/submissions/:id/decision", handlers.DecideSubmission)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("HSE backend (Go) running at http://localhost:%s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
