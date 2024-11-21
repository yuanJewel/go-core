package task

import (
	"github.com/yuanJewel/go-core/db/redis"
	"os"
	"strconv"
)

type Task struct {
	Tag         string      `required:"true" yaml:"tag" env:"task.tag"`
	Concurrency int         `default:"10" yaml:"concurrency" env:"task.concurrency"`
	IsWorker    bool        `default:"false" yaml:"worker" env:"task.worker"`
	Redis       redis.Redis `yaml:"redis"`
	RabbitMq    `yaml:"mq"`
}

type RabbitMq struct {
	Host     string `required:"true" yaml:"host" env:"task.mq.host"`
	Port     string `default:"5672" yaml:"port" env:"task.mq.port"`
	Username string `required:"true" yaml:"username" env:"task.mq.username"`
	Password string `required:"true" yaml:"password" env:"task.mq.password"`
	Queue    string `default:"machinery_tasks" yaml:"queue" env:"task.mq.queue"`
	Exchange string `default:"machinery_exchange" yaml:"exchange" env:"task.mq.exchange"`
}

func maxConcurrency() int {
	concurrency := os.Getenv("TASK_MAX_CONCURRENCY")
	if concurrency != "" {
		concurrencyNumber, err := strconv.Atoi(concurrency)
		if err == nil {
			return concurrencyNumber
		}
	}
	return 5
}
