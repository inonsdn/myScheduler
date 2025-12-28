package servicehandler

import "scheduler/internal/config"

type RabbitMQService struct {
	opts config.Options
}

func NewRabbitMQService(opts config.Options) *RabbitMQService {
	return &RabbitMQService{
		opts: opts,
	}
}

func (r *RabbitMQService) Run() {

}

func (r *RabbitMQService) RegisterRoute() {

}

func (r *RabbitMQService) OnShutdown() {

}
