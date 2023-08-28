package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"log"
)

const TaskSendVerifyEmail = "task:send-verify-email"

type PayloadSendVerifyEmail struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func (processor *RedisQueueProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal TaskSendVerifyEmail: %w", asynq.SkipRetry)
	}

	to := []string{payload.Email}
	var attachments []string

	if err := processor.emailSender.SendEmail("verify", payload, to, attachments); err != nil {
		return err
	}
	log.Printf("type: %v, name: %v, email: %v, token: %v -> processed task.",
		task.Type(),
		payload.Name,
		payload.Email,
		payload.Token,
	)
	return nil
}
