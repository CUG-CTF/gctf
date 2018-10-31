package gctfConfig

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

var GCTF_DEBUG = true
var GCTF_DB_DRIVER = ""
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
	GCTF_DB_DRIVER = os.Getenv("GCTF_DB_DRIVER")
	if GCTF_DB_DRIVER ==""{
		log.Fatalln("You must set env GCTF_DB_DRIVER & GCTF_DB_STRING")
	}

	GCTF_DB_STRING = os.Getenv("GCTF_DB_STRING")
	if GCTF_DB_STRING == "" {
		log.Fatalln("You must set a env GCTF_DB_STRING")

	}
	if GCTF_DEBUG{
		log.Println("db message:",GCTF_DB_STRING)
	}
}
