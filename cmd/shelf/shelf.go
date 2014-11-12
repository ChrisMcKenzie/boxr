package main

import (
	"log"
	"os"

	"github.com/fsouza/go-dockerclient"
)

func main() {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)

	config := docker.CreateContainerOptions{
		Name: "test-app",
		Config: &docker.Config{
			Image:        "node",
			Cmd:          []string{"ls"},
			AttachStdout: true,
			AttachStderr: true,
			WorkingDir:   "/app",
			Volumes: map[string]struct{}{
				"/app": {},
			},
		},
	}

	container, err := client.CreateContainer(config)

	if err != nil {
		log.Fatal(err)
		return
	}

	hostConfig := docker.HostConfig{
		Binds: []string{
			"/opt/go/src/github.com/Secret-Ironman/boxr:/app:rw",
		},
	}

	// fmt.Printf("%#v", container)

	e := client.StartContainer(container.ID, &hostConfig)

	if e != nil {
		log.Fatal(e)
		return
	}

	error := client.Logs(docker.LogsOptions{
		Container:    container.ID,
		OutputStream: os.Stdout,
		ErrorStream:  os.Stdout,
		Follow:       true,
		Stdout:       true,
		Stderr:       true,
	})

	if error != nil {
		log.Fatal(error)
		return
	}
}
