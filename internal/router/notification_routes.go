package router

import (
	"github.com/gin-gonic/gin"
	"notificationService/internal/delivery"
	"notificationService/internal/service"
)

func RegisterNotificationRoutes(r *gin.Engine, svc service.NotificationService) {
	h := delivery.NewNotificationHandler(svc)
	r.POST("/", h.CreateNotification)
	r.POST("/ws/message", h.SendMessageWS)
	r.GET("/", h.GetUserNotifications)
	r.GET("/:id", h.GetNotificationByID)
	r.PATCH("/:id/read", h.MarkNotificationAsRead)
	r.DELETE("/:id", h.DeleteNotification)
}
