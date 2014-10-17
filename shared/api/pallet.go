package api

import (
	"log"
	"time"

	"github.com/Secret-Ironman/boxr/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/Secret-Ironman/boxr/shared/types"
)

// type Pallet struct {
// 	Id   int
// 	Name string
// 	Url  string
// }

func (a *Api) PalletGetOne(c *gin.Context) {
	start := time.Now()
	name := c.Params.ByName("name")

	var pallet types.Pallet

	err := a.db.SelectOne(&pallet, "select * from pallets where Name=?", name)

	if err != nil {
		c.JSON(500, Response{
			Message: err,
			Took:    time.Since(start),
			Success: false,
		})
		return
	}

	c.JSON(200, Response{
		Message: pallet,
		Took:    time.Since(start),
		Success: true,
	})
}

func (a *Api) PalletGetAll(c *gin.Context) {
	start := time.Now()

	var pallets []types.Pallet
	_, err := a.db.Select(&pallets, "select * from pallets order by name")

	if err != nil {
		c.JSON(500, Response{
			Message: err,
			Took:    time.Since(start),
			Success: false,
		})
		return
	}

	c.JSON(200, Response{
		Message: pallets,
		Took:    time.Since(start),
		Success: true,
	})
}

func (a *Api) PalletCreate(c *gin.Context) {
	start := time.Now()
	var pallet types.Pallet

	if !c.Bind(&pallet) {
		log.Fatal("Unable to bind data.")
		c.JSON(400, Response{Message: "Unable to bind data.", Took: time.Since(start)})
		return
	}

	err := a.db.Insert(&pallet)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, Response{Message: err, Took: time.Since(start), Success: false})
		return
	}

	c.JSON(201, Response{
		Message: pallet,
		Took:    time.Since(start),
		Success: true,
	})
}
