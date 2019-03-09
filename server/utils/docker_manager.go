package utils

import (
	. "../model"
	"github.com/fsouza/go-dockerclient"
	"log"
)




type DockerPolling struct {
	Clients  []*docker.Client
	Number  uint64
	Current uint64
}

func NewPollingDockerClient() *DockerPolling {
	r := new(DockerPolling)
	var clients []*docker.Client
	for _, x := range GCTFConfig.GCTF_DOCKERS {
		DockerClient, err := docker.NewClient(x)
		if err != nil {
			log.Println("error to connect " + x + " :" + err.Error())
			continue
		}
		clients = append(clients, DockerClient)
	}
	r.Clients = clients
	r.Number = uint64(len(clients))
	r.Current = 0
	if r.Number == 0 {
		log.Fatal("Can't connect any docker server")
	}
	return r
}
func (polling *DockerPolling) GetDockerClient() *docker.Client {
	client := polling.Clients[polling.Current%polling.Number]
	polling.Current++
	return client
}
