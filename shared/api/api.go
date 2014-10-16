package api

import (
	"time"

	"github.com/Secret-Ironman/boxr/Godeps/_workspace/src/github.com/gin-gonic/gin"
)

type Response struct {
	Message interface{}   `json:"message" binding:"required"`
	Took    time.Duration `json:"took" binding:"required"`
}

func Run(port string) {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/commits/:id", CommitsHookGet)
		api.POST("/pallet", PalletCreate)
	}

	r.Run(port)
}
