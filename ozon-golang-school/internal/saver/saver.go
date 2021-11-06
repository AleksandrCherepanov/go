package saver

import (
	"time"

	"github.com/ozonva/ova-link-api/internal/flusher"
	"github.com/ozonva/ova-link-api/internal/link"
)

type Saver interface {
	Save(entity link.Link)
	Close()
}

type saveWorker struct {
	save  chan link.Link
	close chan bool
}
type timeoutSaver struct {
	entities []link.Link
	flusher  flusher.Flusher
	capacity uint
	ticker   *time.Ticker
	worker   saveWorker
}

func NewTimeOutSaver(capacity uint, flusher flusher.Flusher, savePeriodInSeconds uint) Saver {
	ts := &timeoutSaver{
		entities: make([]link.Link, 0, capacity),
		flusher:  flusher,
		capacity: capacity,
		ticker:   time.NewTicker(time.Second * time.Duration(savePeriodInSeconds)),
		worker: saveWorker{
			save:  make(chan link.Link),
			close: make(chan bool),
		},
	}

	ts.startWorker()
	return ts
}

func (ts *timeoutSaver) Save(entity link.Link) {
	ts.worker.save <- entity
}

func (ts *timeoutSaver) Close() {
	ts.worker.close <- true
	close(ts.worker.save)
	close(ts.worker.close)
}

func (ts *timeoutSaver) addToFlush(entity link.Link) {
	if len(ts.entities) == int(ts.capacity) {
		ts.flush()
	}
	ts.entities = append(ts.entities, entity)
}

func (ts *timeoutSaver) flush() {
	unprocessed := ts.flusher.Flush(ts.entities)
	if len(unprocessed) > 0 {
		for _, entity := range unprocessed {
			ts.addToFlush(entity)
		}
		return
	}
	ts.entities = ts.entities[:0]
}

func (ts *timeoutSaver) startWorker() {
	go func(ts *timeoutSaver) {
	exit:
		for {
			select {
			case <-ts.ticker.C:
				ts.flush()
			case <-ts.worker.close:
				ts.ticker.Stop()
				ts.flush()
				break exit
			case entity := <-ts.worker.save:
				ts.addToFlush(entity)
			}
		}
	}(ts)
}
