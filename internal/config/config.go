package config

import (
	"fmt"
	"os"
)

// Main options for app
type Options struct {
	localDbOptions *LocalDbOptions
	serverOptions  *ServerOptions
	lineOptions    *LineOptions
}

type LocalDbOptions struct {
	runningPort int
	host        string
	port        int
	user        string
	password    string
	dbName      string
}

type OptsFunc func(*Options)

func defaultOptions() Options {
	lineOptions := GetLineOptions()
	serverOptions := GetServerOptions()
	return Options{
		lineOptions:   lineOptions,
		serverOptions: serverOptions,
	}
}

func GetLineOptions() *LineOptions {
	channelSecret := os.Getenv("CHANNEL_SECRET")
	accessToken := os.Getenv("CHANNEL_ACCESS_TOKEN")
	return &LineOptions{
		webhookUrl:    "/webhook",
		port:          3000,
		channelSecret: channelSecret,
		accessToken:   accessToken,
	}
}

func GetServerOptions() *ServerOptions {
	return &ServerOptions{
		host: "localhost",
		port: 8080,
	}
}

func GetLocalDbOptions() *LocalDbOptions {
	return &LocalDbOptions{
		user:     "",
		password: "",
		dbName:   "",
	}
}

func GetOptions(funcs ...OptsFunc) *Options {
	opts := defaultOptions()
	for _, f := range funcs {
		f(&opts)
	}
	return &opts
}

func (l *LocalDbOptions) GetAddress() string {
	return fmt.Sprintf("%s:%d", l.host, l.runningPort)
}

func (l *LocalDbOptions) DataSourceOptions() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		l.user,
		l.password,
		l.host,
		l.port,
		l.dbName,
	)
}

func (o *Options) GetLineOptions() *LineOptions {
	return o.lineOptions
}

func (o *Options) GetServerOptions() *ServerOptions {
	return o.serverOptions
}

func (o *Options) GetLocalDbOptions() *LocalDbOptions {
	return o.localDbOptions
}

type LineOptions struct {
	webhookUrl    string
	port          int
	channelSecret string
	accessToken   string
}

func (l *LineOptions) GetWebhookUrl() string {
	return l.webhookUrl
}

func (l *LineOptions) GetPort() int {
	return l.port
}

func (l *LineOptions) GetChannelSecret() string {
	return l.channelSecret
}

func (l *LineOptions) GetAccessToken() string {
	return l.accessToken
}

type ServerOptions struct {
	host string
	port int
}
