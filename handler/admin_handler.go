package handler

import (
	"api/entity"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AdminHandler serves summary data for the web admin dashboard.
type AdminHandler struct {
	db *gorm.DB
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{db: db}
}

// Stats returns row counts for each managed resource. "users" counts only
// app users (role user), and new_users_* count their registrations.
func (h *AdminHandler) Stats(c *gin.Context) {
	now := time.Now()
	queries := map[string]*gorm.DB{
		"locations":     h.db.Model(&entity.Location{}),
		"categories":    h.db.Model(&entity.Category{}),
		"events":        h.db.Model(&entity.Event{}),
		"trails":        h.db.Model(&entity.Trail{}),
		"reviews":       h.db.Model(&entity.Review{}),
		"users":         h.db.Model(&entity.User{}).Where("role = ?", entity.RoleUser),
		"new_users_7d":  h.db.Model(&entity.User{}).Where("role = ? AND created_at >= ?", entity.RoleUser, now.AddDate(0, 0, -7)),
		"new_users_30d": h.db.Model(&entity.User{}).Where("role = ? AND created_at >= ?", entity.RoleUser, now.AddDate(0, 0, -30)),
	}
	out := make(map[string]int64, len(queries))
	for key, q := range queries {
		var n int64
		if err := q.Count(&n).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		out[key] = n
	}
	c.JSON(http.StatusOK, out)
}

// AdminUser is one registered account as shown in the admin user list.
type AdminUser struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Role          string    `json:"role"`
	AvatarURL     string    `json:"avatar_url"`
	LoginMethod   string    `json:"login_method"` // "google" or "email"
	ReviewCount   int64     `json:"review_count"`
	BookmarkCount int64     `json:"bookmark_count"`
	CreatedAt     time.Time `json:"created_at"`
}

// Users lists every registered account, newest first, with how many reviews
// each has written. Passwords are never selected.
func (h *AdminHandler) Users(c *gin.Context) {
	var users []AdminUser
	err := h.db.Model(&entity.User{}).
		Select(`users.id, users.name, users.email, users.role, users.created_at, users.avatar_url,
			CASE WHEN users.google_id IS NULL THEN 'email' ELSE 'google' END AS login_method,
			(SELECT COUNT(*) FROM reviews WHERE reviews.user_id = users.id) AS review_count,
			(SELECT COUNT(*) FROM bookmarks WHERE bookmarks.user_id = users.id) AS bookmark_count`).
		Order("users.created_at DESC").
		Scan(&users).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Registrations per day for the last 30 days (app users only), oldest first.
	type dayCount struct {
		Day   string `json:"day"`
		Count int64  `json:"count"`
	}
	start := time.Now().AddDate(0, 0, -29)
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.Local)
	byDay := map[string]int64{}
	for _, u := range users {
		if u.Role == entity.RoleUser && !u.CreatedAt.Before(start) {
			byDay[u.CreatedAt.In(time.Local).Format(dateLayout)]++
		}
	}
	daily := make([]dayCount, 0, 30)
	for d := start; len(daily) < 30; d = d.AddDate(0, 0, 1) {
		key := d.Format(dateLayout)
		daily = append(daily, dayCount{Day: key, Count: byDay[key]})
	}

	if users == nil {
		users = make([]AdminUser, 0)
	}
	c.JSON(http.StatusOK, gin.H{"data": users, "daily": daily})
}
