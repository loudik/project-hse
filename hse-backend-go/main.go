package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"hse-backend-go/db"
	"hse-backend-go/handlers"
	"hse-backend-go/middleware"
	"hse-backend-go/utils"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	db.Connect()
	db.ConnectRedis()

	go func() {
		time.Sleep(30 * time.Second)
		utils.CheckExpiringCertificates()
		utils.CheckExpiringAuthorisations()
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			utils.CheckExpiringCertificates()
			utils.CheckExpiringAuthorisations()
		}
	}()
	if err := utils.InitStorage(); err != nil {
		log.Fatalf("Failed to initialize MinIO storage: %v", err)
	}

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
		protected.GET("/notifications", handlers.ListNotifications)
		protected.GET("/notifications/unread-count", handlers.GetUnreadNotificationCount)
		protected.PATCH("/notifications/read-all", handlers.MarkAllNotificationsRead)
		protected.PATCH("/notifications/:id/read", handlers.MarkNotificationRead)

		protected.POST("/organizations", handlers.CreateOrganization)
		protected.GET("/organizations/mine", handlers.GetMyOrganization)

		reviewOnly := protected.Group("/organizations")
		reviewOnly.Use(middleware.RequireRole("Admin", "ANP HSE"))
		{
			reviewOnly.GET("", handlers.ListOrganizations)
			reviewOnly.GET("/:id", handlers.GetOrganization)
			reviewOnly.GET("/:id/history", handlers.GetOrganizationHistory)
			reviewOnly.PATCH("/:id/decision", handlers.DecideOrganization)
		}

		vesselReviewOnly := protected.Group("/vessel-applications")
		vesselReviewOnly.Use(middleware.RequireRole("Admin", "ANP HSE"))
		{
			vesselReviewOnly.GET("/review", handlers.ListVesselApplicationsForReview)
			// vesselReviewOnly.GET("/review/:id", handlers.GetVesselApplicationForReview)
			// vesselReviewOnly.GET("/review/:id/documents", handlers.ListVesselDocumentsForReview)
			vesselReviewOnly.PATCH("/:id/decision", handlers.DecideVesselApplication)
			vesselReviewOnly.PATCH("/:id/assign", handlers.AssignHSEOfficer)
			vesselReviewOnly.GET("/hse-officers", handlers.ListHSEOfficers)
			vesselReviewOnly.POST("/:id/authorisation", handlers.IssueAuthorisationLetter)
			vesselReviewOnly.GET("/extension-requests", handlers.ListExtensionRequests)
			vesselReviewOnly.PATCH("/extension-requests/:id/decision", handlers.DecideExtension)
		}

		dashboardGroup := protected.Group("/dashboard")
		dashboardGroup.Use(middleware.RequireRole("Admin", "ANP HSE"))
		{
			dashboardGroup.GET("/summary", handlers.GetDashboardSummary)
			dashboardGroup.GET("/trend", handlers.GetDashboardTrend)
			dashboardGroup.GET("/avg-approval-time", handlers.GetDashboardAvgApprovalTime)
		}

		protected.POST("/vessel-applications", handlers.CreateVesselApplication)
		protected.GET("/vessel-applications", handlers.ListMyVesselApplications)
		protected.GET("/vessel-applications/:id", handlers.GetVesselApplication)
		protected.PATCH("/vessel-applications/:id", handlers.UpdateVesselApplication)
		protected.POST("/vessel-applications/:id/submit", handlers.SubmitVesselApplication)
		protected.GET("/vessel-applications/lookup", handlers.LookupVesselApplication)
		protected.POST("/vessel-applications/:id/withdraw", handlers.WithdrawVesselApplication)
		protected.POST("/vessel-applications/:id/reopen", handlers.ReopenVesselApplication)
		protected.POST("/vessel-applications/:id/documents", handlers.UploadVesselDocument)
		protected.GET("/vessel-applications/:id/documents", handlers.ListVesselDocuments)
		protected.GET("/vessel-applications/:id/documents/:docId/download", handlers.DownloadVesselDocument)
		protected.GET("/vessel-applications/:id/outcome-document/download", handlers.DownloadVesselOutcomeDocument)
		protected.GET("/vessel-applications/:id/outcome-document/view", handlers.ViewVesselOutcomeDocument)
		protected.GET("/vessel-applications/:id/history", handlers.GetVesselApplicationHistory)
		protected.GET("/vessel-applications/:id/authorisation", handlers.GetAuthorisationLetter)
		protected.GET("/vessel-applications/:id/authorisation/download", handlers.DownloadAuthorisationLetter)
		protected.POST("/vessel-applications/:id/extension", handlers.RequestExtension)
		protected.GET("/vessel-applications/:id/extension-requests", handlers.ListExtensionRequestsForApplication)
		protected.POST("/vessel-applications/substitute", handlers.SubstituteVessel)
		protected.GET("/vessel-applications/review/:id", handlers.GetVesselApplicationForReview)
		protected.GET("/vessel-applications/review/:id/documents", handlers.ListVesselDocumentsForReview)
		protected.GET("/vessel-applications/staff-users", handlers.ListStaffUsers)
		protected.POST("/vessel-applications/:id/review-requests", handlers.RequestStaffReview)
		protected.GET("/vessel-applications/:id/review-requests", handlers.ListReviewRequests)
		protected.POST("/vessel-applications/:id/comments", handlers.AddReviewComment)
		protected.GET("/vessel-applications/:id/comments", handlers.ListReviewComments)
		protected.POST("/vessel-applications/:id/notify-anp-ready", handlers.NotifyANPReady)

		adminOnly := protected.Group("")
		adminOnly.Use(middleware.RequireRole("Admin"))
		{
			adminOnly.GET("/users", handlers.ListUsers)
			adminOnly.GET("/roles", handlers.ListRoles)
			adminOnly.PATCH("/users/:id/role", handlers.UpdateUserRole)
			adminOnly.PATCH("/users/:id/status", handlers.UpdateUserStatus)
			adminOnly.GET("/admin/access-matrix", handlers.GetAccessMatrix)
			adminOnly.PATCH("/admin/access-matrix", handlers.UpdateAccessMatrix)
			adminOnly.GET("/departments", handlers.ListDepartments)
			adminOnly.POST("/departments", handlers.CreateDepartment)
			adminOnly.PATCH("/departments/:id", handlers.UpdateDepartment)
			adminOnly.DELETE("/departments/:id", handlers.DeleteDepartment)
			adminOnly.GET("/departments/:id/members", handlers.ListDepartmentMembers)
			adminOnly.PATCH("/users/:id/department", handlers.UpdateUserDepartment)
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
