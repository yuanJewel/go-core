package task

import (
	"fmt"
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/sirupsen/logrus"
	"github.com/yuanJewel/go-core/db/redis"
	"github.com/yuanJewel/go-core/logger"
)

const StateAborted = "ABORTED"

var (
	MachineryInstance *machinery.Server
)

func InitWork(redis redis.Redis, task Task, taskMap map[string]interface{}) (err error) {
	rabbitmq := task.RabbitMq
	MachineryInstance, err = machinery.NewServer(&config.Config{
		Broker:          fmt.Sprintf("amqp://%s:%s@%s:%s", rabbitmq.Username, rabbitmq.Password, rabbitmq.Host, rabbitmq.Port),
		DefaultQueue:    rabbitmq.Queue,
		ResultBackend:   fmt.Sprintf("redis://%s@%s:%s/%d", redis.Password, redis.Host, redis.Port, redis.Db),
		ResultsExpireIn: 3600,
		Redis: &config.RedisConfig{
			MaxIdle:      redis.PoolSize,
			ReadTimeout:  redis.Timeout,
			WriteTimeout: redis.Timeout,
		},
		AMQP: &config.AMQPConfig{
			Exchange:      rabbitmq.Exchange,
			ExchangeType:  "direct",
			BindingKey:    rabbitmq.Queue,
			PrefetchCount: task.Concurrency,
		},
	})
	if err != nil {
		return
	}

	taskMap["success"] = resultToDb
	taskMap["error"] = errorToDb
	taskMap["finish"] = finishTask

	err = MachineryInstance.RegisterTasks(taskMap)
	if err != nil {
		return
	}
	log.SetInfo(&wrapper{logrus.InfoLevel})
	log.SetDebug(&wrapper{logrus.DebugLevel})
	log.SetError(&wrapper{logrus.ErrorLevel})
	log.SetWarning(&wrapper{logrus.WarnLevel})

	logger.Log.Infof("Complete task system registration !")
	if task.IsWorker {
		worker := MachineryInstance.NewWorker(task.Tag, task.Concurrency)
		logger.Log.Infof("Start one worker !")
		worker.LaunchAsync(make(chan error))
	}
	return
}
