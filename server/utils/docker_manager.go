package utils

import "github.com/fsouza/go-dockerclient"

type DockerManager interface {
	GetDockerClient() *docker.Client
}

type Polling struct {
	Clients []*docker.Client
}
