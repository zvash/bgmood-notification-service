package worker

import (
	"context"
	"github.com/hibiken/asynq"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

// QueueProcessor gets fed by new tasks from redis and processes them
type QueueProcessor interface {
	Start() error
	Run() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}
