package config

import (
	"fmt"
	"os"
	"strconv"
)

func SetPort() OptsFunc {
	return func(opts *Options) {
		o := opts.lineOptions
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
