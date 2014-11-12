package build

import (
	"fmt"
	"testing"

	"github.com/Secret-Ironman/boxr/shared/git"
	"github.com/Secret-Ironman/boxr/shared/parser"
	"github.com/fsouza/go-dockerclient"
)

func TestRun(t *testing.T) {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	b := New(client)

	b.Boxr = &parser.Boxr{
		Box:  "boxr/node",
		Name: "test_pallet",
	}

	b.Repo = &git.Repo{
		Name: "test_pallet",
		Path: "https://github.com/boxr-io/test_pallet.git",
		Dir:  fmt.Sprintf("/var/repos/%s", "test_pallet"),
	}

	b.Run()
}
