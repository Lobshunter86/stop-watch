package main

import (
	"sync"
	"time"
)

const LongDuration = 24 * time.Hour * 365 * 100

type Ticker struct {
	rwLock   *sync.RWMutex
	ticker   *time.Ticker
	duration time.Duration
	stopped  bool
}

func NewTicker(d time.Duration) *Ticker {
	return &Ticker{
		duration: d,
		ticker:   time.NewTicker(LongDuration),
		rwLock:   &sync.RWMutex{},
		stopped:  true,
	}
}

func (t *Ticker) Stop() {
	t.rwLock.Lock()
	defer t.rwLock.Unlock()

	if !t.stopped {
		t.stopped = true
		t.ticker.Reset(LongDuration) // reset interval to a long duration, result like stopped
	}
}

func (t *Ticker) Start() {
	t.rwLock.Lock()
	defer t.rwLock.Unlock()

	if t.stopped {
		t.stopped = false
		t.ticker.Reset(t.duration)
	}
}
