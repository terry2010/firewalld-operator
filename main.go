package main

import (
	"./common"
	"./op"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"log"
	"os"
)

var ok bool
var err error

func main() {
	initConfig()
	r := initRouter()
	err = r.Run(":" + Common.Config.GetString("http.port") + "")

}

func initRouter() *gin.Engine {
	router := gin.Default()
	//router.Use(ginzap.Ginzap(Common.Logger, time.RFC3339, true))
	//router.Use(ginzap.RecoveryWithZap(Common.Logger, true))

	router.NoRoute(Common.Page404)

	router.GET("/firewall/add", op.CrontrollerFirewallRichRuleAdd)

	log.Println(os.Getpid(), "Server-Started")
	log.Println(os.Getpid(), "Server-IP:PORT:", Common.GetServerIP()+":"+Common.Config.GetString("http.port"))

	return router
}

func initConfig() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	var helpDOC = `
Usage: server [OPTIONS] [cmd [arg [arg ...]]]

Examples:
	server --mode=debug --port=7788 
	  

Note:
When no command is given, server starts in debug mode.
Type "--help" or "-h" in interactive mode for information on available commands
and settings.
`

	h := pflag.BoolP("help", "h", false, helpDOC)

	runMode := pflag.StringP("mode", "m", "debug", "server run mode")

	port := pflag.IntP("port", "p", -1, "server listening port[local]")

	pflag.Parse()

	if true == *h {
		log.Println(os.Getpid(), helpDOC)
		os.Exit(107)
		return
	}
	runPath, err := Common.GetCurrentPath()
	if nil != err {
		panic("Get CurrentPath Error:" + runPath + "  |  " + Common.SafeGetError(err))
	}

	Common.Config.SetConfigType("json")

	Common.Config.SetConfigName("config")
	Common.Config.AddConfigPath(runPath + "config/" + *runMode)
	err = Common.Config.ReadInConfig()

	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	log.Println(os.Getpid(), "init config:RunPath:", runPath, " | RunMode:", *runMode)

	Common.Config.Set("system.runPath", runPath)
	Common.Config.Set("system.runMode", *runMode)
	Common.Config.WatchConfig()
	Common.Config.OnConfigChange(func(in fsnotify.Event) {

		log.Println(os.Getpid(), "Config file changed:", in.Name, in.Op.String())
		err = Common.Config.ReadInConfig()

		if err != nil { // Handle errors reading the config file
			log.Println(os.Getpid(), fmt.Errorf("Fatal error config file: %s \n", err))
		} else {
			log.Println(os.Getpid(), "reload config success")
		}
	})
	if 1 > *port {
		panic("need port")
	}

	Common.Config.Set("http.ip", Common.GetServerIP())
	Common.Config.Set("http.port", port)

}
