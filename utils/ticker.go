package utils

import "time"

type TStatus uint8

const (
	RESET TStatus = iota
	STOP
)

type Ticker struct {
	ticker    *time.Ticker
	status    chan TStatus
	resetTime time.Duration
}

type Option func(*Ticker)

func WithResetTime(d time.Duration) Option {
	return func(t *Ticker) {
		t.resetTime = d
	}
}

func NewTicker(opts ...Option) *Ticker {
	defaultTime := 6 * time.Second
	t := &Ticker{
		ticker:    time.NewTicker(defaultTime),
		status:    make(chan TStatus),
		resetTime: defaultTime,
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

func (t *Ticker) Stop() {
	t.status <- STOP
}

func (t *Ticker) Reset() {
	t.status <- RESET
}

func (t *Ticker) Wait() {
	for {
		select {
		case <-t.ticker.C:
			return
		case status := <-t.status:
			if status == STOP {
				t.ticker.Stop()
			} else if status == RESET {
				t.ticker.Reset(t.resetTime)
			}
		}
	}
}
