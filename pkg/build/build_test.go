package build

import (
	"fmt"
	"testing"

	"github.com/Secret-Ironman/boxr/pkg/docker"
	"github.com/Secret-Ironman/boxr/pkg/git"
	"github.com/Secret-Ironman/boxr/pkg/parser"
)

func TestRun(t *testing.T) {
	client := docker.New()
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
