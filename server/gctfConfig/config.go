package gctfConfig

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"log"
	"os"
	"strconv"
)

var (
	GCTF_DEBUG     bool
	GCTF_DB_DRIVER string
	GCTF_DB_STRING string
	GCTF_DOMAIN    string
	//TODO: add docker server manager,else use local docker unix sock
	GCTF_DOCKERS string
	//only in dev
	DockerClient *docker.Client
)

func init() {
	initCreateDir()
	initConnetDocker()
	if i := os.Getenv("GCTF_DEBUG"); i != "" {
		b, e := strconv.ParseBool(i)
		if e != nil {
			log.Fatal("GCTF_DEBUG must be true/false or 0/1")
		}
		GCTF_DEBUG = b
		fmt.Println("You are in product mode")
	} else {
		log.Println("Your are in DEBUG mode!")
	}
	GCTF_DB_DRIVER = os.Getenv("GCTF_DB_DRIVER")
	GCTF_DB_STRING = os.Getenv("GCTF_DB_STRING")
	GCTF_DOMAIN = os.Getenv("GCTF_DOMAIN")
	GCTF_DOCKERS = os.Getenv("GCTF_DOCKERS")
	if GCTF_DB_DRIVER == "" {
		log.Fatalln("You must set env GCTF_DB_DRIVER & GCTF_DB_STRING")
	}
	if GCTF_DB_STRING == "" {
		log.Fatalln("You must set a env GCTF_DB_STRING")
	}
	if GCTF_DOMAIN == "" {
		GCTF_DOMAIN = "localhost"
		if GCTF_DEBUG {
			log.Println("conf.domain: " + GCTF_DOMAIN)
		}
	}
	if GCTF_DEBUG {
		log.Println("conf.db message:", GCTF_DB_STRING)
	}
	if GCTF_DOCKERS == "" {
		log.Println("You are not set DOCKER server,will us local unix sock")
	}
}

func initCreateDir() {
	_ = os.Mkdir("problem", os.ModeDir)
}

func initConnetDocker() {
	var err error
	DockerClient, err = docker.NewClient("unix:///var/run/docker.sock")
	if err != nil {
		log.Fatal("docker server:error connect to local unix sock")
	}
}
