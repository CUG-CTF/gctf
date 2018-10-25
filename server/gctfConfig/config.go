package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

var GCTF_DEBUG = false

func init() {
	if i:=os.Getenv("GCTF_DEBUG");i!=""{
		b,e:=strconv.ParseBool(i)
		if e!=nil{
			log.Fatal("GCTF_DEBUG must be true/false or 0/1")
		}
		GCTF_DEBUG=b
		fmt.Println("You are in product mode")
	} else{
		log.Println("Your are in DEBUG mode!")
	}
}