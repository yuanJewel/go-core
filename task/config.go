package task

import (
	"github.com/yuanJewel/go-core/config"
	"os"
	"strconv"
)

type Task struct {
	Tag         string       `required:"true" yaml:"tag" env:"task.tag"`
	Concurrency int          `default:"10" yaml:"concurrency" env:"task.concurrency"`
	IsWorker    bool         `default:"false" yaml:"worker" env:"task.worker"`
	Redis       config.Redis `yaml:"redis"`
	RabbitMq    `yaml:"mq"`
	// ResultsExpiration is Task result expiration time.
	// After the task ends, the parameters and lock retention time will be recycled according to this time.
	// If the task fails and is blocked, the recycling time will increase exponentially.
	// The storage time is evaluated by the redis service pressure.
	// It is generally set to 1 minute and is recommended to be no less than 15 seconds.
	ResultsExpiration int `default:"60" yaml:"results_expiration" env:"task.results_expiration"`
	// LockExpiration is Task atomic protection lock expiration time.
	// Set according to the estimated maximum time for the task. You can set a longer time to enhance protection.
	LockExpiration int `default:"18000" yaml:"lock_expiration" env:"task.lock_expiration"`
	// VarExpiration is Task parameter expiration time.
	VarExpiration int `default:"300" yaml:"var_expiration" env:"task.var_expiration"`
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
