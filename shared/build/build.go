package build

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/Secret-Ironman/boxr/shared/dockerfile"
	"github.com/Secret-Ironman/boxr/shared/git"
	"github.com/Secret-Ironman/boxr/shared/parser"
	"github.com/Secret-Ironman/boxr/shared/utils"
	"github.com/fsouza/go-dockerclient"
)

var log = utils.Logger()

func New(dockerClient *docker.Client) *Builder {
	return &Builder{
		dockerClient: dockerClient,
	}
}

type Builder struct {
	Boxr *parser.Boxr
	Repo *git.Repo

	dockerClient *docker.Client
	// image created by the builder
	image *docker.Image
	// container to run
	container *docker.Container
	// service containers
	services *[]docker.Container
}

func (b *Builder) Run() error {
	// 1.) setup image and service containers
	b.setup()
	// 2.) run image
	return nil
}

func (b *Builder) setup() {

	// 1.) create tar stream
	inputbuf := bytes.NewBuffer(nil)
	tr := tar.NewWriter(inputbuf)
	b.writeDockerfile(tr)
	// 2.) add source files to tar
	b.writeSourceDir(tr)
	tr.Close()

	// 3.) build docker images
	// make sure the image isn't empty. this would be bad
	if len(b.Boxr.Box) == 0 {
		log.Error("Fatal Error, No Docker Image specified")
		// return fmt.Errorf("Error: missing Docker image")
	}
}

func (b *Builder) writeDockerfile(tr *tar.Writer) {
	t := time.Now()
	df := dockerfile.New(b.Boxr.Box)
	df.WriteEnv("BOXR", "true")
	tr.WriteHeader(&tar.Header{Name: "Dockerfile", Size: int64(len(df.Bytes())), ModTime: t, AccessTime: t, ChangeTime: t})
	fmt.Printf("%v", df)
	tr.Write(df.Bytes())
}

func (b *Builder) writeSourceDir(tr *tar.Writer) error {
	walkFn := func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsDir() {
			return nil
		}
		// Because of scoping we can reference the external root_directory variable
		new_path := path[len(b.Repo.Dir):]
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

	return filepath.Walk(b.Repo.Dir, walkFn)
}
