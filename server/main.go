package main

import (
	. "./model"
	. "./utils"
	"./v1"
	"encoding/json"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func init() {
	//config
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

	//database
	var err error
	// go-xorm is used to create database engine
	// engine, err := xorm.NewEngine(driverName, dataSourceName)
	// data
	GctfDataManage, err = xorm.NewEngine(GCTFConfig.GCTF_DB_DRIVER, GCTFConfig.GCTF_DB_STRING)
	// All table name have a gctf_ prefix

	// prefix，前缀
	tbMapper := core.NewPrefixMapper(core.GonicMapper{}, "gctf_")

	// fix problem_I_D to problem_id
	GctfDataManage.SetColumnMapper(core.GonicMapper{})
	GctfDataManage.SetTableMapper(tbMapper)

	// Ping is test the database is alive
	err = GctfDataManage.Ping()
	if err != nil {
		log.Fatal("database connect error: ", err.Error())
	}
	if GCTFConfig.GCTF_DEBUG {
		//GctfDataManage.ShowSQL(true)
		//GctfDataManage.Logger().SetLevel(core.LOG_DEBUG)
	}
	// this is create lots of tables?
	err = GctfDataManage.CreateTables(User{}, Problems{}, UserProblems{}, Hints{}, Teams{})
	// GctfDataManage.DropTables("gctf_user","gctf_problems","gctf_user_problems","gctf_hints","gctf_tag","gctf_teams")
	checkerr(err)
}

func checkerr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func readConf() {
	confFile, err := os.Open("conf.json")
	if GCTFConfig == nil {
		GCTFConfig = new(GCTFConfigStruct)
	}
	if err != nil {
		log.Fatal("error to open conf file " + err.Error())
	}
	conf, err := ioutil.ReadAll(confFile)
	if err != nil {
		log.Fatal("error to read conf file " + err.Error())
	}
	err = json.Unmarshal(conf, GCTFConfig)
	if err != nil {
		log.Fatal("error to read json conf" + err.Error())
	}
}
func createDir() {
	//TODO:add problem's dir option
	_ = os.Mkdir("problem", os.ModeDir)
}

func connetDocker() {
	//TODO:select docker manager type
	GCTFDockerManager = NewPollingDockerClient()
}

//Todo:更多单元测试
func main() {
	gCTFRoute := gin.Default()
	gCTFRoute.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	v1.ConfigRoute(gCTFRoute.Group("/v1"))
	err := gCTFRoute.Run(":8081")
	if err != nil {
		log.Fatal(err)
	}
}
