package core

import (
	"context"
	"github.com/rendau/my-otus/task8/sender/internal/domain/entities"
	"github.com/rendau/my-otus/task8/sender/internal/interfaces"
	"sync"
)

// Core - is core of logic
type Core struct {
	log           interfaces.Logger
	mq            interfaces.Mq
	ctx           context.Context
	ctxCancelFunc context.CancelFunc
	wg            *sync.WaitGroup
}

// CreateCore - creates core instance
func CreateCore(log interfaces.Logger, mq interfaces.Mq) *Core {
	return &Core{
		log: log,
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
	var event *entities.Event

	defer ucs.ctxCancelFunc()
	defer ucs.wg.Done()

	for {
		event, err = ucs.mq.GetEvent(ucs.ctx)
		if err != nil {
			ucs.log.Errorw("Fail to get event from mq", "error", err.Error())
			return
		}

		if event == nil {
			return
		}

		ucs.log.Infof("Message from mq: %+v", event)
	}
}

// Stop - stops core
func (ucs Core) Stop() {
	ucs.ctxCancelFunc()
	ucs.wg.Wait()
}
