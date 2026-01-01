package config

import (
	"os"
	"strconv"
)

func SetPort() OptsFunc {
	return func(opts *Options) {
		o := opts.serverOptions
		port := os.Getenv("PORT")
		port_int, err := strconv.Atoi(port)
		if err != nil {
			o.port = 8080
		}
		o.port = int(port_int)
	}
}
