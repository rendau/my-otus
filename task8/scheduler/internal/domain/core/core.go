package core

import (
	"context"
	"github.com/rendau/my-otus/task8/scheduler/internal/domain/entities"
	"github.com/rendau/my-otus/task8/scheduler/internal/interfaces"
	"sync"
	"time"
)

// Core - is core of logic
type Core struct {
	log           interfaces.Logger
	stg           interfaces.Storage
	mq            interfaces.Mq
	ctx           context.Context
	ctxCancelFunc context.CancelFunc
	wg            *sync.WaitGroup
}

// CreateCore - creates core instance
func CreateCore(log interfaces.Logger, stg interfaces.Storage, mq interfaces.Mq) *Core {
	return &Core{
		log: log,
		stg: stg,
		mq:  mq,
		wg:  &sync.WaitGroup{},
	}
}

// Run - returns list of event
func (ucs *Core) Run() {
	ucs.ctx, ucs.ctxCancelFunc = context.WithCancel(context.Background())

	ucs.wg.Add(1)
	go ucs.run()
}

// Done - returns ctx done-channel
func (ucs Core) Done() <-chan struct{} {
	return ucs.ctx.Done()
}

func (ucs *Core) run() {
	var err error
	var t1 time.Time
	var t2 time.Time
	var events []*entities.Event
	var event *entities.Event

	defer ucs.ctxCancelFunc()
	defer ucs.wg.Done()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ucs.ctx.Done():
			return
		case <-ticker.C:
			t1 = time.Now().Truncate(24 * time.Hour)
			t2 = time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour)

			// read from storage
			events = nil
			events, err = ucs.stg.EventList(context.Background(), &entities.EventListFilter{
				StartTimeGt: &t1,
				StartTimeLt: &t2,
			})
			if err != nil {
				ucs.log.Errorw("Fail to get event list", "error", err.Error())
				return
			}

			// publish to queue
			if len(events) > 0 {
				for _, event = range events {
					err = ucs.mq.PublishEventNotification(event)
					if err != nil {
						ucs.log.Errorw("Fail to publish notification to mq", "error", err.Error())
						return
					}
				}
			}
		}
	}
}

// Stop - stops core
func (ucs Core) Stop() {
	ucs.ctxCancelFunc()
	ucs.wg.Wait()
}
