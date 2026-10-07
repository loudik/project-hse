package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"hse-backend-go/db"
)

// GET /api/dashboard/summary - ANP HSE/Admin only
// Counts of vessel applications by status, plus organizations pending review.
func GetDashboardSummary(c *gin.Context) {
	statusCounts := map[string]int{}
	rows, err := db.DB.Query(`SELECT status, COUNT(*) FROM vessel_applications GROUP BY status`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch summary"})
		return
	}
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err == nil {
			statusCounts[status] = count
		}
	}
	rows.Close()

	var pendingOrganizations int
	_ = db.DB.QueryRow(`SELECT COUNT(*) FROM organizations WHERE status = 'Pending'`).Scan(&pendingOrganizations)

	c.JSON(http.StatusOK, gin.H{
		"vesselApplications":   statusCounts,
		"pendingOrganizations": pendingOrganizations,
	})
}

// GET /api/dashboard/trend?months=6 - ANP HSE/Admin only
// Vessel applications submitted per month, for a bar/line chart.
func GetDashboardTrend(c *gin.Context) {
	months := 6
	if m := c.Query("months"); m != "" {
		if parsed, err := strconv.Atoi(m); err == nil && parsed > 0 && parsed <= 24 {
			months = parsed
		}
	}

	rows, err := db.DB.Query(`
		SELECT DATE_FORMAT(submitted_at, '%Y-%m') AS month, COUNT(*) AS count
		FROM vessel_applications
		WHERE submitted_at IS NOT NULL
		  AND submitted_at >= DATE_SUB(CURDATE(), INTERVAL ? MONTH)
		GROUP BY month
		ORDER BY month ASC
	`, months)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch trend"})
		return
	}
	defer rows.Close()

	type point struct {
		Month string `json:"month"`
		Count int    `json:"count"`
	}
	list := []point{}
	for rows.Next() {
		var p point
		if err := rows.Scan(&p.Month, &p.Count); err == nil {
			list = append(list, p)
		}
	}
	c.JSON(http.StatusOK, list)
}

// GET /api/dashboard/avg-approval-time - ANP HSE/Admin only
// Average time from submission to decision, across decided applications.
func GetDashboardAvgApprovalTime(c *gin.Context) {
	var avgHours *float64
	err := db.DB.QueryRow(`
		SELECT AVG(TIMESTAMPDIFF(HOUR, submitted_at, decided_at))
		FROM vessel_applications
		WHERE submitted_at IS NOT NULL AND decided_at IS NOT NULL
	`).Scan(&avgHours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to compute average approval time"})
		return
	}

	if avgHours == nil {
		c.JSON(http.StatusOK, gin.H{"avgHours": nil, "avgDays": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"avgHours": *avgHours,
		"avgDays":  *avgHours / 24,
	})
}
