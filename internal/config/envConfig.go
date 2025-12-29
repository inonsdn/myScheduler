package config

import "os"

func SetPort() OptsFunc {
	return func(opts *Options) {
		o := opts.serverOptions
		env := os.Getenv("ENV")
		if env == "prod" {
			o.port = 443
		} else {
			o.port = 8080
		}
	}
}
