package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"log"
)

const TaskSendResetPasswordEmail = "task:send-password-reset-email"

type PayloadSendResetPasswordEmail struct {
	AppName string `json:"app_name"`
	Email   string `json:"email"`
	Token   string `json:"token"`
}

func (processor *RedisQueueProcessor) ProcessTaskSendResetPasswordEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendResetPasswordEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal TaskSendResetPasswordEmail: %w", asynq.SkipRetry)
	}

	to := []string{payload.Email}
	var attachments []string

	if err := processor.emailSender.SendEmail("reset-password", payload, to, attachments); err != nil {
		return err
	}
	log.Printf("type: %v, app name: %v, email: %v, token: %v -> processed task.",
		task.Type(),
		payload.AppName,
		payload.Email,
		payload.Token,
	)
	return nil
}
