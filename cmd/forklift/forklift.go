package main

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	data "github.com/Secret-Ironman/boxr/shared/db"
	"github.com/Secret-Ironman/boxr/shared/git"
	"github.com/Secret-Ironman/boxr/shared/parser"
	"github.com/Secret-Ironman/boxr/shared/types"
	"github.com/Secret-Ironman/boxr/shared/utils"
	"github.com/fsouza/go-dockerclient"
)

func main() {
	log := utils.Logger()

	db, err := data.New("boxr.db")
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

		// parse boxr.yml
		boxr, err := parser.ParseBoxrFile(fmt.Sprintf("%s/boxr.yml", repo.Dir))

		if err != nil {
			log.Error(err.Error())
		}

		log.Debug("%#v", boxr)

		pallet.Status = "building"

		_, err = db.Update(&pallet)
		if err != nil {
			log.Error(err.Error())
		}

		// Build docker containers
		endpoint := "unix:///var/run/docker.sock"
		client, _ := docker.NewClient(endpoint)

		t := time.Now()
		inputbuf, outputbuf := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
		file := []byte{}
		box := fmt.Sprintf("FROM %v\n", boxr.Box)
		file = append(file, []byte(box)...)
		boxr.Build.RunSteps(&file, "cd /tmp")
		file = append(file, []byte("RUN mkdir -p /opt/app && cp -a /tmp/node_modules /opt/app/")...)
		volume := "WORKDIR /opt/app \nADD . /opt/app\n"
		file = append(file, []byte(volume)...)
		file = append(file, []byte("EXPOSE 3000\n")...)

		file = append(file, []byte(fmt.Sprintf("CMD %s\n", boxr.Run))...)

		log.Debug(string(file))
		tr := tar.NewWriter(inputbuf)
		tr.WriteHeader(&tar.Header{Name: "Dockerfile", Size: int64(len(file)), ModTime: t, AccessTime: t, ChangeTime: t})
		tr.Write(file)
		walkFn := func(path string, info os.FileInfo, err error) error {
			if info.Mode().IsDir() {
				return nil
			}
			// Because of scoping we can reference the external root_directory variable
			new_path := path[len(repo.Dir):]
			if len(new_path) == 0 {
				return nil
			}
			fr, err := os.Open(path)
			if err != nil {
				return err
			}
			defer fr.Close()

			if h, err := tar.FileInfoHeader(info, new_path); err != nil {
				log.Error(err.Error())
			} else {
				h.Name = new_path
				if err = tr.WriteHeader(h); err != nil {
					log.Error(err.Error())
				}
			}
			if _, err := io.Copy(tr, fr); err != nil {
				log.Error(err.Error())
			}
			return nil
		}

		if err = filepath.Walk(repo.Dir, walkFn); err != nil {
			log.Error(err.Error())
		}
		tr.Close()

		opts := docker.BuildImageOptions{
			Name:         boxr.Name,
			InputStream:  inputbuf,
			OutputStream: outputbuf,
		}
		if err := client.BuildImage(opts); err != nil {
			log.Fatal(err)
		}

		container := docker.CreateContainerOptions{
			Name: boxr.Name,
			Config: &docker.Config{
				Image:        boxr.Name,
				AttachStdout: true,
				AttachStderr: true,
				Tty:          true,
			},
		}

		var c *docker.Container
		c, err = client.CreateContainer(container)
		if err != nil {
			log.Error(err.Error())
		}

		hostconfig := docker.HostConfig{
			/*PortBindings: map[docker.Port][]docker.PortBinding{
				"3000": []docker.PortBinding{
					docker.PortBinding{
						HostIp:   "0.0.0.0",
						HostPort: "3000",
					},
				},
			},*/
			PublishAllPorts: true,
		}

		err = client.StartContainer(c.ID, &hostconfig)
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
