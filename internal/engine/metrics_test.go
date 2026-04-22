package engine

import (
	"sync"
	"testing"
)

func TestMetricsConcurrentEvents(t *testing.T) {
	m := NewMetrics(2048)

	var wg sync.WaitGroup
	for i := 0; i < 500; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			m.RecordCorrect()
		}()
		go func() {
			defer wg.Done()
			m.RecordIncorrect()
		}()
	}
	wg.Wait()

	// flush ensures the collector has processed all buffered events
	m.Flush()

	correct, incorrect := m.Counts()
	m.Close()
	if correct != 500 {
		t.Errorf("got correct=%d, want 500", correct)
	}
	if incorrect != 500 {
		t.Errorf("got incorrect=%d, want 500", incorrect)
	}
}

func TestMetricsElapsed(t *testing.T) {
	m := NewMetrics(64)
	defer m.Close()

	elapsed := m.Elapsed()
	if elapsed <= 0 {
		t.Error("elapsed should be positive after construction")
	}
}

func TestMetricsInitialCountsZero(t *testing.T) {
	m := NewMetrics(64)
	defer m.Close()

	correct, incorrect := m.Counts()
	if correct != 0 || incorrect != 0 {
		t.Errorf("fresh metrics should be zero, got correct=%d incorrect=%d", correct, incorrect)
	}
}
