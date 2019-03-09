package config

import (
	"encoding/json"
	"github.com/fsouza/go-dockerclient"
	"io/ioutil"
	"log"
	"os"
)

type gctfConfig struct {
	GCTF_DEBUG     bool `json:"debug"`
	GCTF_DB_DRIVER string `json:"database_type"`
	GCTF_DB_STRING string `json:"database_address"`
	GCTF_DOMAIN    string  `json:"domain_name"`
	//TODO: add docker server manager,else use local docker unix sock
	GCTF_DOCKERS []string `json:"docker_servers"`
}
var (
	GCTFConfig gctfConfig
	//TODO: add docker server manager,else use local docker unix sock
	//only in dev2
	DockerClient *docker.Client
)

func init() {
	readConf()
	createDir()
	connetDocker()
	if GCTFConfig.GCTF_DEBUG {
		log.Println("You are in product mode")
	} else {
		log.Println("Your are in DEBUG mode!")
	}
	if GCTFConfig.GCTF_DB_DRIVER == "" {
		log.Fatalln("You must set env GCTF_DB_DRIVER & GCTF_DB_STRING")
	}
	if GCTFConfig.GCTF_DB_STRING == "" {
		log.Fatalln("You must set a env GCTF_DB_STRING")
	}
	if GCTFConfig.GCTF_DOMAIN == "" {
		GCTFConfig.GCTF_DOMAIN = "localhost"
		if GCTFConfig.GCTF_DEBUG {
			log.Println("conf.domain: " + GCTFConfig.GCTF_DOMAIN)
		}
	}
	if GCTFConfig.GCTF_DEBUG {
		log.Println("conf.db message:", GCTFConfig.GCTF_DB_STRING)
	}
	if GCTFConfig.GCTF_DOCKERS == nil {
		log.Println("You are not set DOCKER server,will us local unix sock")
	}
}

func readConf(){
	confFile,err:=os.Open("conf.json")
	if err!=nil{
		log.Fatal("error to open conf file"+err.Error())
	}
	conf,err:=ioutil.ReadAll(confFile)
	if err!=nil{
		log.Fatal("error to read conf file"+err.Error())
	}
	err=json.Unmarshal(conf,&GCTFConfig)
	if err!=nil{
		log.Fatal("error to read json conf"+err.Error())
	}
}
func createDir() {
	_ = os.Mkdir("problem", os.ModeDir)
}

func connetDocker() {
	var err error
	DockerClient, err = docker.NewClient("unix:///var/run/docker.sock")
	if err != nil {
		log.Fatal("docker server:error connect to local unix sock")
	}
}
