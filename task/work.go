package task

import (
	"fmt"
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/sirupsen/logrus"
	"github.com/yuanJewel/go-core/db/redis"
	"github.com/yuanJewel/go-core/logger"
	"net"
	"sync"
	"time"
)

const StateAborted = "ABORTED"

var (
	machineryInstance *machinery.Server
	redisInstance     *redis.Store
	lockExpiration    time.Duration
	varExpiration     time.Duration
	finishExpiration  time.Duration
	stepToJob         sync.Map
)

func InitWork(task Task, taskMap map[string]interface{}, f FinishInterface) (err error) {
	lockExpiration = time.Duration(task.LockExpiration) * time.Second
	varExpiration = time.Duration(task.VarExpiration) * time.Second
	finishExpiration = time.Duration(task.ResultsExpiration) * time.Second
	stepToJob = sync.Map{}
	machineryInstance, err = machinery.NewServer(&config.Config{
		Broker:          fmt.Sprintf("amqp://%s:%s@%s:%s", task.RabbitMq.Username, task.RabbitMq.Password, task.RabbitMq.Host, task.RabbitMq.Port),
		DefaultQueue:    task.RabbitMq.Queue,
		ResultBackend:   fmt.Sprintf("redis://%s@%s:%s/%d", task.Redis.Password, task.Redis.Host, task.Redis.Port, task.Redis.Db),
		ResultsExpireIn: task.ResultsExpiration,
		Redis: &config.RedisConfig{
			MaxIdle:      task.Redis.PoolSize,
			ReadTimeout:  task.Redis.Timeout,
			WriteTimeout: task.Redis.Timeout,
		},
		AMQP: &config.AMQPConfig{
			Exchange:      task.RabbitMq.Exchange,
			ExchangeType:  "direct",
			BindingKey:    task.RabbitMq.Queue,
			PrefetchCount: task.Concurrency,
		},
	})
	if err != nil {
		return
	}

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(task.RabbitMq.Host, task.RabbitMq.Port), 3*time.Second)
	if err != nil || conn == nil {
		return fmt.Errorf("cannot connect task.RabbitMq(%s:%s), error: %v", task.RabbitMq.Host, task.RabbitMq.Port, err)
	}
	_ = conn.Close()
	conn, err = net.DialTimeout("tcp", net.JoinHostPort(task.Redis.Host, task.Redis.Port), 3*time.Second)
	if err != nil || conn == nil {
		return fmt.Errorf("cannot connect redis(%s:%s), error: %v", task.Redis.Host, task.Redis.Port, err)
	}
	_ = conn.Close()

	redisInstance, err = redis.GetRedisInstance(&task.Redis)
	if err != nil {
		return
	}

	taskMap["success"] = resultToDb
	taskMap["error"] = errorToDb
	taskMap["finish"] = finishTask

	err = machineryInstance.RegisterTasks(taskMap)
	if err != nil {
		return
	}
	log.SetInfo(&wrapper{logrus.InfoLevel})
	log.SetDebug(&wrapper{logrus.DebugLevel})
	log.SetError(&wrapper{logrus.ErrorLevel})
	log.SetWarning(&wrapper{logrus.WarnLevel})

	logger.Log.Infof("Complete task system registration !")
	if task.IsWorker {
		worker := machineryInstance.NewWorker(task.Tag, task.Concurrency)
		logger.Log.Infof("Start one worker !")
		worker.LaunchAsync(make(chan error))
	}
	finishObject = f
	return
}
