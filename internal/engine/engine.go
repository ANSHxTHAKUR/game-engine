package engine

import (
	"context"
	"errors"
	"sync"

	"github.com/ansh-singh/game-engine/internal/model"
)

type GameEngine interface {
	Submit(response model.UserResponse) error
	Result() <-chan model.GameResult
	GetMetrics() *Metrics
	Shutdown()
}

type gameEngine struct {
	responses     chan model.UserResponse
	resultCh      chan model.GameResult
	metrics       *Metrics
	cancel        context.CancelFunc
	ctx           context.Context
	declareWinner sync.Once
	shutdownOnce  sync.Once
	winnerID      string
	shutdownCh    chan struct{}
	wg            sync.WaitGroup
}

func New(workerCount, bufferSize int) GameEngine {
	ctx, cancel := context.WithCancel(context.Background())
	ge := &gameEngine{
		responses:  make(chan model.UserResponse, bufferSize),
		resultCh:   make(chan model.GameResult, 1),
		metrics:    NewMetrics(bufferSize),
		cancel:     cancel,
		ctx:        ctx,
		shutdownCh: make(chan struct{}),
	}
	ge.startWorkers(workerCount)
	return ge
}

func (ge *gameEngine) Submit(resp model.UserResponse) error {
	select {
	case <-ge.shutdownCh:
		return errors.New("engine is shut down")
	default:
	}
	select {
	case ge.responses <- resp:
		return nil
	case <-ge.shutdownCh:
		return errors.New("engine is shut down")
	}
}

func (ge *gameEngine) Result() <-chan model.GameResult {
	return ge.resultCh
}

func (ge *gameEngine) GetMetrics() *Metrics {
	return ge.metrics
}

func (ge *gameEngine) Shutdown() {
	ge.shutdownOnce.Do(func() {
		close(ge.shutdownCh)
		ge.cancel()
		ge.wg.Wait()
		close(ge.responses)
		ge.metrics.Close()
	})
}

func (ge *gameEngine) startWorkers(count int) {
	for i := 0; i < count; i++ {
		ge.wg.Add(1)
		go ge.worker()
	}
}
