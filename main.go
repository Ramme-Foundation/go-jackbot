package main

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go-jackbot/internal/controller"
	"go-jackbot/prisma/db"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func main() {
	_ = godotenv.Load()
	logger, _ := zap.NewProduction()

	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		logger.With(zap.Error(err)).Fatal("failed to connect to prisma")
	}

	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			logger.With(zap.Error(err)).Fatal("failed to disconnect from prisma")
		}
	}()

	controller := controller.NewController(client)

	r := gin.New()
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))

	// err := r.SetTrustedProxies(nil)
	/* if err != nil {
		panic(err)
	}*/
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/command", controller.CommandRoute())
	err := r.Run()
	if err != nil {
		logger.With(zap.Error(err)).Fatal("failed to start server")
	}
}
