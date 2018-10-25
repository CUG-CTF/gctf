package gctfConfig

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

var GCTF_DEBUG = true
var GCTF_DB_NAME = ""
var GCTF_DB_STRING = ""

func init() {
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
	GCTF_DB_NAME = os.Getenv("GCTF_DB_NAME")
	if GCTF_DB_NAME ==""{
		log.Fatalln("You must set env GCTF_DB_NAME")
	}

	GCTF_DB_STRING = os.Getenv("GCTF_DB_STRING")
	if GCTF_DB_STRING == "" {
		log.Fatalln("You must set a env GCTF_DB_STRING")

	}
	if GCTF_DEBUG{
		log.Println("db message:",GCTF_DB_STRING)
	}
}
