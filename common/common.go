package Common

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func InitConfig() (err error) {
	runPath, err := GetCurrentPath()

	Config.SetConfigType("json")

	Config.SetConfigName("config")
	Config.AddConfigPath(runPath)
	err = Config.ReadInConfig()

	if err != nil { // Handle errors reading the config file
		log.Println(runPath)
		return errors.New(fmt.Sprintf("Fatal error config file: %s \n", err))
	}

	Config.WatchConfig()
	Config.OnConfigChange(func(in fsnotify.Event) {
		log.Println(os.Getpid(), "Config file changed:", in.Name, in.Op.String())
		err = Config.ReadInConfig()

		if err != nil { // Handle errors reading the config file
			log.Println(os.Getpid(), fmt.Errorf("Fatal error config file: %s \n", err))
		} else {
			log.Println(os.Getpid(), "reload config success")
		}
	})
	return
}

func GetCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	//if err != nil {
	//	return "", err
	//}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}

	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}

	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\".`)
	}
	return string(path[0 : i+1]), nil
}

func GetServerIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Println(os.Getpid(), err)
		return ""
	}
	for _, value := range addrs {
		if ipnet, ok := value.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.IsGlobalUnicast() && !ipnet.IP.IsInterfaceLocalMulticast() {
			if ipnet.IP.To4() != nil {
				//log.Println(os.Getpid(),ipnet.IP.String())
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func Page404(c *gin.Context) {
	//返回404状态码
	c.JSON(http.StatusNotFound, gin.H{
		"code": http.StatusNotFound,
		"msg":  "404, page not exists!",
		"data": "",
	})
}

func FastAtoi(num string) int {
	ret, _ := strconv.Atoi(num)
	return ret
}

func FastJsonMarshal(_json interface{}) string {
	str, _ := Json.MarshalToString(_json)
	return str
}

func Md5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func SafeGetError(err error) string {
	if nil == err {
		return ""
	} else {
		return err.Error()
	}
}
