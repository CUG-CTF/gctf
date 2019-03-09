package utils

import (
	"github.com/fsouza/go-dockerclient"
	"log"
	. "../config"
)

type DockerManager interface {
	GetDockerClient() *docker.Client
}

type DockerPolling struct {
	clients []*docker.Client
	number  uint64
	current uint64
}

func NewPollingDockerClient() *DockerPolling{
	r:=new(DockerPolling)
	var clients []*docker.Client
	for _, x := range GCTFConfig.GCTF_DOCKERS {
		DockerClient, err := docker.NewClient(x)
		if err != nil {
			log.Println("error to connect " + x + " :" + err.Error())
			continue
		}
		clients= append(clients, DockerClient)
	}
	r.clients=clients
	r.number=uint64(len(clients))
	r.current=0
	if r.number==0{
		log.Fatal("Can't connect any docker server")
	}
	return r
}
func (polling *DockerPolling) GetDockerClient() *docker.Client {
	client := polling.clients[polling.current%polling.number]
	polling.current++
	return client
}
