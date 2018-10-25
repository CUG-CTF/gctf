package main

import (
	"fmt"
	"os"
)

const GCTF_Debug=false

func init()  {
	if i:=os.Getenv("GCTF_DEBUG");i!=""{
		GCTF_Debug=strconv.A
	}
	fmt.Println("123")
	fmt.Println(debug)
}
func main() {

}
