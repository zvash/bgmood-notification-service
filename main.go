package main

import (
	"github.com/hibiken/asynq"
	"github.com/zvash/bgmood-notification-service/internal/mail"
	"github.com/zvash/bgmood-notification-service/internal/util"
	"github.com/zvash/bgmood-notification-service/internal/worker"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	redisOpt := createRedisClientOption(config)
	runTaskProcessor(config, redisOpt)
}

func createRedisClientOption(config util.Config) asynq.RedisClientOpt {
	return asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
}

func runTaskProcessor(config util.Config, redisOpt asynq.RedisClientOpt) {
	emailSender := mail.NewGeneralEmailSender(config)
	asynqConfig := worker.AsynqConfig{
		Concurrency: config.WorkerCount,
	}
	taskProcessor := worker.NewRedisQueueProcessor(redisOpt, emailSender, asynqConfig)
	log.Println("start task processor")
	err := taskProcessor.Run()
	if err != nil {
		log.Fatal("failed to start task processor")
	}
}
