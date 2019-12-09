package rmq

import (
	"context"
	"encoding/json"
	"github.com/rendau/my-otus/task8/sender/internal/domain/entities"
	"github.com/streadway/amqp"
)

const queueName = "event_notify"

// Rmq - is type for rabbit-mq client
type Rmq struct {
	con         *amqp.Connection
	ch          *amqp.Channel
	consumeChan <-chan amqp.Delivery
}

// NewRmq - creates new Rmq instance
func NewRmq(dsn string) (*Rmq, error) {
	var err error

	res := &Rmq{}

	res.con, err = amqp.Dial(dsn)
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

	res.consumeChan, err = res.ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetEvent - gets event from mq
func (r *Rmq) GetEvent(ctx context.Context) (*entities.Event, error) {
	var err error
	var delivery amqp.Delivery
	var event entities.Event

	select {
	case <-ctx.Done():
		break
	case delivery = <-r.consumeChan:
		err = json.Unmarshal(delivery.Body, &event)
		if err != nil {
			return nil, err
		}
		return &event, nil
	}

	return nil, nil
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
