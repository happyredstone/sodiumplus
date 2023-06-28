package web

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateServer() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	server := gin.New()

	server.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	server.Use(gin.Recovery())

	server.Use(func(c *gin.Context) {
		dir, err := os.Getwd()

		if err != nil {
			c.String(500, err.Error())
			return
		}

		path := dir + c.Request.URL.Path

		stat, err := os.Stat(path)

		if err != nil {
			c.String(404, err.Error())
			return
		}

		if !stat.IsDir() {
			c.File(path)
		}

		content, err := CreatePage(path)

		if err != nil {
			c.String(500, err.Error())
			return
		}

		c.Header("Content-Type", "text/html")

		c.String(200, content)
	})

	return server
}

func RunServer(server *gin.Engine, address ...string) error {
	addr := "0.0.0.0:4000"

	if len(address) > 0 {
		addr = address[0]
	}

	fmt.Println("Server running on " + addr + "...")

	return server.Run(addr)
}
