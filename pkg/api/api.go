package api

import (
	"github.com/letsu/evping/pkg/hosts"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/letsu/evping/pkg/value"
)

func Router(host *hosts.HostsCsv) {
	r := gin.New()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "UserID", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/api")

	api.Use(hostMiddleware(host))

	api.GET("/allhosts", value.GetHosts)
	api.GET("/dataofhost", value.DataOfHost)
	api.POST("/host", value.AddHost)
	api.DELETE("/host", value.DeleteHost)

	r.Run(":8081")
}
func hostMiddleware(host *hosts.HostsCsv) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("hostKey", host)
		c.Next()
	}
}
