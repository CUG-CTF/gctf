package gctfConfig

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

var (
	GCTF_DEBUG     bool
	GCTF_DB_DRIVER string
	GCTF_DB_STRING string
	GCTF_DOMAIN string
)

func init() {
	initCreateDir()
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
	GCTF_DOMAIN    = os.Getenv("GCTF_DOMAIN")
	if GCTF_DB_DRIVER == "" {
		log.Fatalln("You must set env GCTF_DB_DRIVER & GCTF_DB_STRING")
	}

	if GCTF_DB_STRING == "" {
		log.Fatalln("You must set a env GCTF_DB_STRING")
	}
	if GCTF_DOMAIN==""{
		GCTF_DOMAIN="localhost"
		if GCTF_DEBUG{
			log.Println("conf.domain: "+GCTF_DOMAIN)
		}
	}
	if GCTF_DEBUG {
		log.Println("conf.db message:", GCTF_DB_STRING)
	}
}

func initCreateDir() {
	_ = os.Mkdir("problem", os.ModeDir)
}
