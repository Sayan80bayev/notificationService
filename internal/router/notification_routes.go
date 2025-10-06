package delivery

import (
	"github.com/gin-gonic/gin"
	"notificationService/internal/delivery"
	"notificationService/internal/service"
)

func RegisterNotificationRoutes(r *gin.RouterGroup, svc service.NotificationService) {
	h := delivery.NewNotificationHandler(svc)

	r.POST("/", h.CreateNotification)
	r.GET("/", h.GetUserNotifications)
	r.GET("/:id", h.GetNotificationByID)
	r.PATCH("/:id/read", h.MarkNotificationAsRead)
	r.DELETE("/:id", h.DeleteNotification)
}
