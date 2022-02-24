/*
	This worker will watch the redis queue:`QUEUE_KEY_SEND_EAMIL` and send email when a new task is added to the queue.
*/
package workers

import (
	"context"
	"os"
	"os/signal"

	"github.com/labstack/gommon/log"
	"github.com/thoas/bokchoy"
)

var QUEUE_KEY_SEND_EAMIL = "tasks.send_email"
var QUEUE_ADDR_SEND_EMAIL = "localhost:6379"

var engine *bokchoy.Bokchoy

func SendEmailWorker() {
	// Consumer worker
	var ctx context.Context

	// Queue gets or creates a new queue
	engine.Queue(QUEUE_KEY_SEND_EAMIL).HandleFunc(func(r *bokchoy.Request) error {
		log.Info("Receive request", r)
		log.Info("Payload:", r.Task.Payload)

		// TODO: Send email

		return nil
	}, bokchoy.WithConcurrency(5)).OnCompleteFunc(
		func(r *bokchoy.Request) error {
			log.Info("Request", r, "has been completed")
			return nil
		},
	)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			log.Info("Received signal, gracefully stopping")
			engine.Stop(ctx)
		}
	}()

	engine.Run(ctx)
}

func PushRequest(payload interface{}) {
	// Producer
	var ctx context.Context
	task, err := engine.Queue(QUEUE_KEY_SEND_EAMIL).Publish(ctx, payload)
	if err != nil {
		log.Fatal(err)
	}
	log.Info(task, "has been published")

}

func init() {
	var err error
	var ctx context.Context
	engine, err = bokchoy.New(ctx, bokchoy.Config{
		Broker: bokchoy.BrokerConfig{
			Type: "redis",
			Redis: bokchoy.RedisConfig{
				Type: "client",
				Client: bokchoy.RedisClientConfig{
					Addr: QUEUE_ADDR_SEND_EMAIL,
				},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
