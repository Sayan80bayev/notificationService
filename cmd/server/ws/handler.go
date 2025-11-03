package ws

import (
	"github.com/Sayan80bayev/go-project/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleWebSocket(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id missing in context"})
		return
	}
	uid := userID.(uuid.UUID)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upgrade websocket"})
		return
	}

	client := NewClient(uid, conn)
	Register(client)

	go client.WritePump()
	go client.ReadPump()
}

func SetupWebSocketRoutes(r *gin.Engine, JWKSUrl string) {
	r.GET("/ws", middleware.AuthMiddleware(JWKSUrl), HandleWebSocket)
}
