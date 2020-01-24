package rmq

import (
	"encoding/json"
	"github.com/rendau/my-otus/task8/scheduler/internal/domain/entities"
	"github.com/streadway/amqp"
	"time"
)

const (
	connectionWaitTimout = 30 * time.Second
	queueName            = "event_notify"
)

// Rmq - is type for rabbit-mq client
type Rmq struct {
	con *amqp.Connection
	ch  *amqp.Channel
}

// NewRmq - creates new Rmq instance
func NewRmq(dsn string) (*Rmq, error) {
	var err error

	res := &Rmq{}

	res.con, err = res.connectionWait(dsn, connectionWaitTimout)
	if err != nil {
		return nil, err
	}

	res.ch, err = res.con.Channel()
	if err != nil {
		return nil, err
	}

	_, err = res.ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Rmq) connectionWait(dsn string, timeout time.Duration) (*amqp.Connection, error) {
	var err error
	var res *amqp.Connection
	var tryCnt int

	deadline := time.Now().Add(timeout)

	for tryCnt < 2 || time.Now().Before(deadline) {
		tryCnt++
		res, err = amqp.Dial(dsn)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}

	return res, err
}

// PublishEventNotification - publishes event to mq
func (r *Rmq) PublishEventNotification(event *entities.Event) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = r.ch.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        eventBytes,
		},
	)

	return err
}

// Stop - stops mq
func (r *Rmq) Stop() error {
	err := r.ch.Close()
	if err != nil {
		return err
	}
	err = r.con.Close()
	if err != nil {
		return err
	}
	return nil
}
