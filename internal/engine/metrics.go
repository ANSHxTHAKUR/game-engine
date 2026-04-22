package engine

import (
	"sync"
	"time"
)

type metricEvent struct {
	kind    int
	flushCh chan struct{}
}

const (
	eventCorrect   = 0
	eventIncorrect = 1
	eventFlush     = 2
)

type Metrics struct {
	events    chan metricEvent
	startTime time.Time
	closeOnce sync.Once

	mu        sync.RWMutex
	correct   int64
	incorrect int64
}

func NewMetrics(bufferSize int) *Metrics {
	m := &Metrics{
		events:    make(chan metricEvent, bufferSize),
		startTime: time.Now(),
	}
	go m.collector()
	return m
}

func (m *Metrics) collector() {
	for ev := range m.events {
		switch ev.kind {
		case eventCorrect:
			m.mu.Lock()
			m.correct++
			m.mu.Unlock()
		case eventIncorrect:
			m.mu.Lock()
			m.incorrect++
			m.mu.Unlock()
		case eventFlush:
			close(ev.flushCh)
		}
	}
}

func (m *Metrics) RecordCorrect() {
	m.events <- metricEvent{kind: eventCorrect}
}

func (m *Metrics) RecordIncorrect() {
	m.events <- metricEvent{kind: eventIncorrect}
}

func (m *Metrics) Flush() {
	done := make(chan struct{})
	m.events <- metricEvent{kind: eventFlush, flushCh: done}
	<-done
}

func (m *Metrics) Counts() (correct, incorrect int64) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.correct, m.incorrect
}

func (m *Metrics) Elapsed() time.Duration {
	return time.Since(m.startTime)
}

func (m *Metrics) Close() {
	m.closeOnce.Do(func() {
		close(m.events)
	})
}
