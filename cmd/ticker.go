package main

import (
	"sync"
	"time"
)

const LongDuration = 24 * time.Hour * 365 * 100

type Ticker struct {
	C        <-chan time.Time
	rwLock   *sync.RWMutex
	ticker   *time.Ticker
	duration time.Duration
	stopped  bool
}

func NewTicker(d time.Duration) *Ticker {
	t := time.NewTicker(LongDuration)

	return &Ticker{
		C:        t.C,
		duration: d,
		ticker:   t,
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
