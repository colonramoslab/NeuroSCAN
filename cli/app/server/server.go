package server

import (
	"os"

	"neuroscan/app/router"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Run() {
	godotenv.Load(".env")
	port := os.Getenv("PORT")
	mode := os.Getenv("GIN_MODE")

	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	if port == "" {
		port = "8080"
	}

	r := router.Router()
	r.ForwardedByClientIP = true
	err := r.SetTrustedProxies([]string{"127.0.0.1"})

	if err != nil {
		log.Fatal("Error setting trusted proxies: " + err.Error())
	}

	err = r.Run(":" + port)

	if err != nil {
		log.Fatal("Error starting server: " + err.Error())
	}
}
