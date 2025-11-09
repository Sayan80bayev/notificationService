package main

import (
	"context"
	"github.com/Sayan80bayev/go-project/pkg/logging"
	"github.com/gin-gonic/gin"
	"notificationService/cmd/server/ws"
	"notificationService/internal/bootstrap"
	"notificationService/internal/router"
)

func main() {
	ctn, err := bootstrap.Init()
	if err != nil {
		panic(err)
	}

	log := logging.GetLogger()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(logging.Middleware)

	ctx, cancel := context.WithCancel(context.Background())
	go ctn.Consumer.Start(ctx)
	defer func() {
		cancel()
		ctn.Consumer.Close()
	}()

	router.RegisterNotificationRoutes(r, ctn.NotificationService)
	ws.SetupWebSocketRoutes(r, ctn.JWKSUrl)
	log.Info("server is running on port " + ctn.Config.Port)
	err = r.Run(":" + ctn.Config.Port)
	if err != nil {
		log.Fatal("can't start server")
		return
	}

}
