package config

import (
	"fmt"
	"os"
	"strconv"
)

func SetHost() OptsFunc {
	return func(opts *Options) {
		o := opts.serverOptions
		mode := os.Getenv("RUN_MODE")
		if mode == "prod" {
			o.host = "0.0.0.0"
		} else {
			o.host = "localhost"
		}
	}
}

func SetPort() OptsFunc {
	return func(opts *Options) {
		o := opts.serverOptions
		port := os.Getenv("PORT")
		port_int, err := strconv.Atoi(port)
		if err != nil {
			fmt.Println("Error when set port: ", err)
			o.port = 8081
			return
		}
		o.port = int(port_int)
	}
}
