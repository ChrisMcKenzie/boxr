package api

import "github.com/Secret-Ironman/boxr/Godeps/_workspace/src/github.com/gin-gonic/gin"

func Run(port string) {
	router := gin.Default()

	router.GET("/commits/:id", CommitsHookGet)
	router.Run(port)
}
