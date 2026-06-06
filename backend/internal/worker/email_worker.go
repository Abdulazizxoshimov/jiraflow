package worker

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/email"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/rabbitmq"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

const (
	emailMaxRetries = 3
	emailRetryDelay = 5 * time.Second
)

// EmailJob is the message payload published to QueueEmailSend.
type EmailJob struct {
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
	Attempt int      `json:"attempt"`
}

// EmailWorker consumes email jobs from RabbitMQ and delivers them with retry.
type EmailWorker struct {
	mq     *rabbitmq.Client
	sender email.Sender
	log    logger.Logger
}

func NewEmailWorker(mq *rabbitmq.Client, sender email.Sender, log logger.Logger) *EmailWorker {
	return &EmailWorker{mq: mq, sender: sender, log: log}
}

// Run starts consuming from QueueEmailSend. Blocking — run in a goroutine.
func (w *EmailWorker) Run(ctx context.Context) {
	if err := w.mq.DeclareQueue(rabbitmq.QueueEmailSend); err != nil {
		w.log.Error(ctx, "email_worker: declare queue failed", logger.SafeString("err", err.Error()))
		return
	}

	msgs, err := w.mq.Consume(rabbitmq.QueueEmailSend)
	if err != nil {
		w.log.Error(ctx, "email_worker: consume failed", logger.SafeString("err", err.Error()))
		return
	}

	w.log.Info(ctx, "email_worker: started")
	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-msgs:
			if !ok {
				return
			}
			w.handle(ctx, d)
		}
	}
}

func (w *EmailWorker) handle(ctx context.Context, d amqp.Delivery) {
	var job EmailJob
	if err := json.Unmarshal(d.Body, &job); err != nil {
		w.log.Error(ctx, "email_worker: unmarshal failed", logger.SafeString("err", err.Error()))
		_ = d.Ack(false)
		return
	}

	err := w.sender.SendRaw(ctx, job.To, job.Subject, job.Body)
	if err == nil {
		_ = d.Ack(false)
		w.log.Info(ctx, "email_worker: sent",
			logger.String("to", job.To[0]),
			logger.String("subject", job.Subject),
		)
		return
	}

	w.log.Warn(ctx, "email_worker: send failed",
		logger.SafeString("err", err.Error()),
		logger.String("attempt", string(rune('0'+job.Attempt))),
	)
	_ = d.Ack(false)

	if job.Attempt < emailMaxRetries {
		job.Attempt++
		time.AfterFunc(emailRetryDelay, func() {
			_ = w.mq.Publish(context.Background(), rabbitmq.QueueEmailSend, job)
		})
	} else {
		w.log.Error(ctx, "email_worker: max retries reached, dropping message",
			logger.String("to", job.To[0]),
			logger.String("subject", job.Subject),
		)
	}
}
