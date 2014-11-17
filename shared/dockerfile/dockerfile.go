package dockerfile

import (
	"bytes"
	"fmt"
)

type Dockerfile struct {
	bytes.Buffer
}

func New(from string) *Dockerfile {
	d := Dockerfile{}
	d.WriteFrom(from)
	return &d
}

func (d *Dockerfile) WriteFrom(from string) {
	d.WriteString(fmt.Sprintf("FROM %s\n", from))
}

func (d *Dockerfile) WriteRun(run string) {
	d.WriteString(fmt.Sprintf("RUN %s\n", run))
}

func (d *Dockerfile) WriteAdd(from string, to string) {
	d.WriteString(fmt.Sprintf("ADD %s %s\n", from, to))
}

func (d *Dockerfile) WriteWorkdir(path string) {
	d.WriteString(fmt.Sprintf("WORKDIR %s\n", path))
}

func (d *Dockerfile) WriteUser(user string) {
	d.WriteString(fmt.Sprintf("USER %s\n", user))
}

func (d *Dockerfile) WriteEnv(key string, value string) {
	d.WriteString(fmt.Sprintf("ENV %s %s\n", key, value))
}

func (d *Dockerfile) WriteExpose(port string) {
	d.WriteString(fmt.Sprintf("EXPOSE %s\n", port))
}

func (d *Dockerfile) WriteEntrypoint(path string) {
	d.WriteString(fmt.Sprintf("ENTRYPOINT %s\n", path))
}
