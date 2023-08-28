package worker

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zvash/bgmood-notification-service/internal/mail"
)

type RedisQueueProcessor struct {
	server      *asynq.Server
	emailSender mail.EmailSender
}

type AsynqConfig struct {
	Concurrency int
}

func NewRedisQueueProcessor(redisOptions asynq.RedisClientOpt, emailSender mail.EmailSender, asynqConfig AsynqConfig) QueueProcessor {
	server := asynq.NewServer(
		redisOptions,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
			Concurrency: asynqConfig.Concurrency,
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				fmt.Printf("error: %v, type: %v, payload: %v -> process task failed",
					err,
					task.Type(),
					task.Payload(),
				)
			}),
		},
	)
	return &RedisQueueProcessor{
		server:      server,
		emailSender: emailSender,
	}
}

func (processor *RedisQueueProcessor) createMux() *asynq.ServeMux {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)
	mux.HandleFunc(TaskSendResetPasswordEmail, processor.ProcessTaskSendResetPasswordEmail)
	return mux
}

func (processor *RedisQueueProcessor) Start() error {
	mux := processor.createMux()
	return processor.server.Start(mux)
}

func (processor *RedisQueueProcessor) Run() error {
	mux := processor.createMux()
	return processor.server.Run(mux)
}
