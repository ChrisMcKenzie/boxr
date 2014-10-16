package api

import (
	"time"

	"github.com/Secret-Ironman/boxr/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/Secret-Ironman/boxr/shared/types"
)

func PalletCreate(c *gin.Context) {
	start := time.Now()
	var body types.Pallet

	if !c.Bind(&body) {
		c.JSON(400, Response{Message: "unknown_error", Took: time.Since(start)})
		return
	}
	c.JSON(201, Response{
		Message: body,
		Took:    time.Since(start),
	})
}
