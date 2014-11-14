package main

import (
	"fmt"
	"os"

	"github.com/Secret-Ironman/boxr/shared/build"
	data "github.com/Secret-Ironman/boxr/shared/db"
	"github.com/Secret-Ironman/boxr/shared/docker"
	"github.com/Secret-Ironman/boxr/shared/git"
	"github.com/Secret-Ironman/boxr/shared/parser"
	"github.com/Secret-Ironman/boxr/shared/types"
	"github.com/Secret-Ironman/boxr/shared/utils"
)

func main() {
	log := utils.Logger()

	db, err := data.New("boxr.db")
	client := docker.New()
	if err != nil {
		log.Critical(err.Error())
	}
	// repo := git.Repo{
	// 	Name: "boxr-io/test_pallet",
	// 	Path: "https://github.com/boxr-io/test_pallet.git",
	// }

	var pallets []types.Pallet
	_, err = db.Select(&pallets, "select * from pallets order by name")

	for _, pallet := range pallets {
		log.Info("Building %s \n", pallet.Name)
		repo := git.Repo{
			Name: pallet.Name,
			Path: pallet.Url,
			Dir:  fmt.Sprintf("/var/repos/%s", pallet.Name),
		}

		pallet.Status = "retrieving"

		_, err = db.Update(&pallet)
		if err != nil {
			log.Error(err.Error())
		}

		if _, err := os.Stat(repo.Dir); os.IsNotExist(err) {
			log.Info("Cloning Repo from %s", repo.Path)
			repo.Clone()
		} else {
			log.Info("Pulling repo from remote %s", repo.Path)
			repo.Pull()
		}

		builder := build.New(client)
		boxr, err := parser.ParseBoxrFile(fmt.Sprintf("%s/boxr.yml", repo.Dir))

		if err != nil {
			log.Error(err.Error())
		}

		builder.Repo = &repo
		builder.Boxr = boxr

		pallet.Status = "building"

		_, err = db.Update(&pallet)
		if err != nil {
			log.Error(err.Error())
		}

		err = builder.Run()

		if err != nil {
			log.Error(err.Error())
		}

		log.Debug("Box Running")
		pallet.Status = "running"

		_, err = db.Update(&pallet)
		if err != nil {
			log.Error(err.Error())
		}
	}
}
