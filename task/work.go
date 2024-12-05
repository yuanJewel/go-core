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
	rabbitmq := task.RabbitMq
	redisConf := task.Redis
	lockExpiration = time.Duration(task.LockExpiration) * time.Second
	varExpiration = time.Duration(task.VarExpiration) * time.Second
	finishExpiration = time.Duration(task.ResultsExpiration) * time.Second
	stepToJob = sync.Map{}
	machineryInstance, err = machinery.NewServer(&config.Config{
		Broker:          fmt.Sprintf("amqp://%s:%s@%s:%s", rabbitmq.Username, rabbitmq.Password, rabbitmq.Host, rabbitmq.Port),
		DefaultQueue:    rabbitmq.Queue,
		ResultBackend:   fmt.Sprintf("redis://%s@%s:%s/%d", redisConf.Password, redisConf.Host, redisConf.Port, redisConf.Db),
		ResultsExpireIn: task.ResultsExpiration,
		Redis: &config.RedisConfig{
			MaxIdle:      redisConf.PoolSize,
			ReadTimeout:  redisConf.Timeout,
			WriteTimeout: redisConf.Timeout,
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

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(rabbitmq.Host, rabbitmq.Port), 3*time.Second)
	if err != nil || conn == nil {
		return fmt.Errorf("cannot connect rabbitmq(%s:%s), error: %v", rabbitmq.Host, rabbitmq.Port, err)
	}
	_ = conn.Close()
	conn, err = net.DialTimeout("tcp", net.JoinHostPort(redisConf.Host, redisConf.Port), 3*time.Second)
	if err != nil || conn == nil {
		return fmt.Errorf("cannot connect redis(%s:%s), error: %v", redisConf.Host, redisConf.Port, err)
	}
	_ = conn.Close()

	redisInstance, err = redis.GetRedisInstance(&redisConf)
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
