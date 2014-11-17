package api

import (
	"fmt"
	"os"
	"time"

	"github.com/Secret-Ironman/boxr/pkg/build"
	"github.com/Secret-Ironman/boxr/pkg/docker"
	"github.com/Secret-Ironman/boxr/pkg/git"
	"github.com/Secret-Ironman/boxr/pkg/parser"
	"github.com/Secret-Ironman/boxr/pkg/types"
	"github.com/gin-gonic/gin"
)

// var log = utils.Logger()

func (a *Api) PalletGetOne(c *gin.Context) {
	start := time.Now()
	name := c.Params.ByName("name")

	var pallet types.Pallet

	err := a.db.SelectOne(&pallet, "select * from pallets where Name=?", name)

	if err != nil {
		c.JSON(500, Response{
			Message: err,
			Took:    time.Since(start).String(),
			Success: false,
		})
		return
	}

	c.JSON(200, Response{
		Message: pallet,
		Took:    time.Since(start).String(),
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
			Took:    time.Since(start).String(),
			Success: false,
		})
		return
	}

	c.JSON(200, Response{
		Message: pallets,
		Took:    time.Since(start).String(),
		Success: true,
	})
}

func (a *Api) PalletCreate(c *gin.Context) {
	start := time.Now()
	var pallet types.Pallet

	if !c.Bind(&pallet) {
		log.Fatal("Unable to bind data.")
		c.JSON(400, Response{Message: "Unable to bind data.", Took: time.Since(start).String()})
		return
	}

	err := a.db.Insert(&pallet)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, Response{Message: err, Took: time.Since(start).String(), Success: false})
		return
	}

	c.JSON(201, Response{
		Message: pallet,
		Took:    time.Since(start).String(),
		Success: true,
	})
}

func (a *Api) PalletBuild(c *gin.Context) {
	start := time.Now()
	name := c.Params.ByName("name")

	var pallet types.Pallet

	err := a.db.SelectOne(&pallet, "select * from pallets where Name=?", name)

	if err != nil {
		c.JSON(500, Response{
			Message: err,
			Took:    time.Since(start).String(),
			Success: false,
		})
		return
	}

	// data, _ := json.Marshal(pallet)
	// resp, err := http.Post("http://localhost:3001/pallet", "application/json", bytes.NewBuffer(data))

	// if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	log.Info(resp)
	// }

	// enqueue here
	go func() {
		client := docker.New()
		if err != nil {
			log.Critical(err.Error())
		}

		log.Info("Building %s \n", pallet.Name)
		repo := git.Repo{
			Name: pallet.Name,
			Path: pallet.Url,
			Dir:  fmt.Sprintf("/var/repos/%s", pallet.Name),
		}

		pallet.Status = "retrieving"

		_, err = a.db.Update(&pallet)
		if err != nil {
			log.Error(err.Error())
		}

		if _, err := os.Stat(repo.Dir); os.IsNotExist(err) {
			log.Info("Cloning Repo from %s", repo.Path)
			repo.Clone()
		} else {
			log.Info("Pulling repo from remote %s", repo.Path)
			log.Info(repo.Pull())
		}

		builder := build.New(client)
		boxr, err := parser.ParseBoxrFile(fmt.Sprintf("%s/boxr.yml", repo.Dir))

		if err != nil {
			log.Error(err.Error())
		}

		builder.Repo = &repo
		builder.Boxr = boxr

		pallet.Status = "building"

		_, err = a.db.Update(&pallet)
		if err != nil {
			log.Error(err.Error())
		}

		err = builder.Run()

		if err != nil {
			log.Error(err.Error())
		}

		log.Debug("Box Running")
		pallet.Status = "running"

		_, err = a.db.Update(&pallet)
		if err != nil {
			log.Error(err.Error())
		}
	}()
}
