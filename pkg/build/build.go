package build

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/Secret-Ironman/boxr/pkg/docker"
	"github.com/Secret-Ironman/boxr/pkg/dockerfile"
	"github.com/Secret-Ironman/boxr/pkg/git"
	"github.com/Secret-Ironman/boxr/pkg/parser"
	"github.com/Secret-Ironman/boxr/pkg/utils"
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
	container *docker.Run
	// service containers
	services []*docker.Container
}

func (b *Builder) Run() error {
	// 1.) setup image and service containers
	if err := b.setup(); err != nil {
		return err
	}
	// 2.) run image
	if err := b.run(); err != nil {
		return err
	}
	return nil
}

func (b *Builder) setup() error {

	// 1.) create tar stream
	inputbuf := bytes.NewBuffer(nil)
	tr := tar.NewWriter(inputbuf)
	b.writeDockerfile(tr)
	// 2.) add source files to tar
	b.writeSourceDir(tr)
	tr.Close()

	// 3.) build docker service images
	for _, service := range b.Boxr.Services {

		// Parse the name of the Docker image
		// And then construct a fully qualified image name
		cname := service

		// Get the image info
		img, err := b.dockerClient.Images.Inspect(cname)
		if err != nil {
			// Get the image if it doesn't exist
			if err := b.dockerClient.Images.Pull(cname); err != nil {
				return fmt.Errorf("Error: Unable to pull image %s", cname)
			}

			img, err = b.dockerClient.Images.Inspect(cname)
			if err != nil {
				return fmt.Errorf("Error: Invalid or unknown image %s", cname)
			}
		}

		// debugging
		log.Info("starting service container %s", cname)

		// Run the contianer
		run, err := b.dockerClient.Containers.RunDaemonPorts(cname, img.Config.ExposedPorts)
		if err != nil {
			return err
		}

		// Get the container info
		info, err := b.dockerClient.Containers.Inspect(run.ID)
		if err != nil {
			// on error kill the container since it hasn't yet been
			// added to the array and would therefore not get
			// removed in the defer statement.
			b.dockerClient.Containers.Stop(run.ID, 10)
			b.dockerClient.Containers.Remove(run.ID)
			return err
		}

		// Add the running service to the list
		b.services = append(b.services, info)
	}

	// make sure the image isn't empty. this would be bad
	if len(b.Boxr.Box) == 0 {
		log.Error("Fatal Error, No Docker Image specified")
		// return fmt.Errorf("Error: missing Docker image")
	}

	b.dockerClient.Images.Build(b.Boxr.Name, inputbuf)

	// debugging
	log.Info("building app in %s", b.Repo.Dir)

	// get the image details
	var err error
	b.image, err = b.dockerClient.Images.Inspect(b.Boxr.Name)
	if err != nil {
		// if we have problems with the image make sure
		// we remove it before we exit
		log.Error("failed to verify build image %s", b.Boxr.Name)
		return err
	}

	return nil
}

func (b *Builder) run() error {
	// create and run the container
	conf := docker.Config{
		Image:        b.image.ID,
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
	}

	// configure if Docker should run in privileged mode
	host := docker.HostConfig{}
	host.PortBindings = make(map[docker.Port][]docker.PortBinding)

	// for port, _ := range b.Boxr.Ports {
	// 	host.PortBindings[port] = []PortBinding{{HostIp: "127.0.0.1", HostPort: ""}}
	// }
	//
	host.PortBindings["4000"] = []docker.PortBinding{{HostIp: "0.0.0.0", HostPort: "4000"}}

	log.Notice("starting build %s", b.Boxr.Name)

	// create the container from the image
	run, err := b.dockerClient.Containers.Create(&conf)
	if err != nil {
		return err
	}

	// cache instance of docker.Run
	b.container = run

	// attach to the container
	// go func() {
	// 	b.dockerClient.Containers.Attach(run.ID, &writer{os.Stdout, 0})
	// }()

	// start the container
	if err := b.dockerClient.Containers.Start(run.ID, &host); err != nil {
		return err
	}

	// wait for the container to stop
	// wait, err := b.dockerClient.Containers.Wait(run.ID)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (b *Builder) writeDockerfile(tr *tar.Writer) {
	t := time.Now()
	df := dockerfile.New(b.Boxr.Box)
	df.WriteEnv("BOXR", "true")
	// TODO: make dynamic port selection
	df.WriteEnv("PORT", "4000")

	df.WriteExpose("4000")

	for _, step := range b.Boxr.Build {
		df.WriteRun(fmt.Sprintf("export PORT=4000; %s; %s", b.Repo.Dir, step))
	}

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
