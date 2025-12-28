package config

import "fmt"

type Options struct {
	runningPort int
	host        string
	port        int
	user        string
	password    string
	dbName      string
}

type LocalDbOptions struct {
	user     string
	password string
	dbName   string
}

type OptsFunc func(*Options)

func defaultOptions() Options {
	return Options{
		runningPort: 8083,
		host:        "localhost",
		port:        3306,
		user:        "root",
		password:    "rootPass",
		dbName:      "schedulerdb",
	}
}

func GetOptions(funcs ...OptsFunc) *Options {
	opts := defaultOptions()
	for _, f := range funcs {
		f(&opts)
	}
	return &opts
}

func (o *Options) GetAddress() string {
	return fmt.Sprintf("%s:%d", o.host, o.runningPort)
}

func (o *Options) GetLocalDbOptions() *LocalDbOptions {
	return &LocalDbOptions{
		user:     o.user,
		password: o.password,
		dbName:   o.dbName,
	}
}

func (o *Options) DataSourceConfig() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		o.user,
		o.password,
		o.host,
		o.port,
		o.dbName,
	)
}
