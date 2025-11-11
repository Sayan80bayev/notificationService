package delivery

import (
	"net/http"
	"notificationService/cmd/server/ws"
	"notificationService/internal/model"
	"notificationService/internal/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	svc service.NotificationService
}

func NewNotificationHandler(svc service.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

// CreateNotification godoc
func (h *NotificationHandler) CreateNotification(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id type"})
		return
	}

	var req struct {
		Title   string `json:"title" binding:"required"`
		Message string `json:"message" binding:"required"`
		Type    string `json:"type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	n := &model.Notification{
		UserID:  userID,
		Message: req.Message,
	}

	created, err := h.svc.CreateNotification(c, n)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GetUserNotifications godoc
func (h *NotificationHandler) GetUserNotifications(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id type"})
		return
	}

	// Optional pagination
	limit := 20
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := parsePositiveInt(l); err == nil {
			limit = parsed
		}
	}
	if o := c.Query("offset"); o != "" {
		if parsed, err := parsePositiveInt(o); err == nil {
			offset = parsed
		}
	}

	notifications, err := h.svc.GetNotificationsByUser(c, userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notifications)
}

// GetNotificationByID godoc
func (h *NotificationHandler) GetNotificationByID(c *gin.Context) {
	idStr := c.Param("id")
	nID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification id"})
		return
	}

	notification, err := h.svc.GetNotificationByID(c, nID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if notification == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
		return
	}
	c.JSON(http.StatusOK, notification)
}

// MarkNotificationAsRead godoc
func (h *NotificationHandler) MarkNotificationAsRead(c *gin.Context) {
	idStr := c.Param("id")
	nID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification id"})
		return
	}

	if err := h.svc.MarkNotificationAsRead(c, nID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "marked as read", "read_at": time.Now().UTC()})
}

// DeleteNotification godoc
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	idStr := c.Param("id")
	nID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification id"})
		return
	}

	if err := h.svc.DeleteNotification(c, nID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// local helper
func parsePositiveInt(val string) (int, error) {
	parsed, err := strconv.Atoi(val)
	if err != nil || parsed < 0 {
		return 0, err
	}
	return parsed, nil
}

func (h *NotificationHandler) SendMessageWS(c *gin.Context) {
	var req struct {
		UserUUID string `json:"user_uuid"`
		Message  string `json:"message"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userUUID, err := uuid.Parse(req.UserUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}

	ws.SendNotification(userUUID, req.Message)

	c.JSON(http.StatusOK, gin.H{
		"status":  "sent",
		"message": req.Message,
		"user":    userUUID.String(),
	})
}
