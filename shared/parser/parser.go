package parser

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Boxr struct {
	// Box specifies the name of the
	// Docker container to use for
	// this machine
	Box string
	// Name of the service; this name
	// will be used to register the
	// service with service discovery
	Name string
	// Version of the service to
	// build; used to determine if the
	// whole cluster needs to be rebuilt.
	Version string
	// Services specifies a list of external
	// services that should be built to support
	// the service.
	Services []string
	// Build phase of the service specifies
	// what is necessary for the service to
	// run.
	Build Phase `yaml:"build,omitempty"`
	// Test phase specifies what to do to
	// test the service.
	Test Phase `yaml:"test,omitempty"`
	// Deploy phase specifies what to do to
	// deploy the service completely on its
	// docker container.
	Run string `yaml:"run,omitempty"`
}

type Phase []string

func ParseBoxr(data string) (*Boxr, error) {
	boxr := Boxr{}

	err := yaml.Unmarshal([]byte(data), &boxr)
	return &boxr, err
}

func ParseBoxrFile(file string) (*Boxr, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return ParseBoxr(string(data))
}

func (p *Phase) RunSteps(file *[]byte, prefix string) {
	for _, step := range *p {
		step = fmt.Sprintf("RUN %s; %s\n", prefix, step)
		*file = append(*file, []byte(step)...)
	}
}
