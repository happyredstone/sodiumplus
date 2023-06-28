package helpers

import (
	"os"
	"strconv"

	"github.com/NoSadBeHappy/SodiumPlus/builder/internal/web"
)

func Serve() error {
	server := web.CreateServer()

	port := 4000

	envPort, exist := os.LookupEnv("PORT")

	if exist {
		port, _ = strconv.Atoi(envPort)
	}

	return web.RunServer(server, "0.0.0.0:"+strconv.Itoa(port))
}
