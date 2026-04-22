package engine

import (
	"fmt"

	"github.com/ansh-singh/game-engine/internal/model"
)

func (ge *gameEngine) worker() {
	defer ge.wg.Done()
	for {
		select {
		case <-ge.ctx.Done():
			ge.drainRemaining()
			return
		case resp, ok := <-ge.responses:
			if !ok {
				return
			}
			ge.evaluate(resp)
		}
	}
}

func (ge *gameEngine) evaluate(resp model.UserResponse) {
	if resp.IsCorrect {
		ge.metrics.RecordCorrect()
		ge.declareWinner.Do(func() {
			ge.winnerID = resp.UserID
			ge.metrics.Flush()
			c, ic := ge.metrics.Counts()
			elapsed := ge.metrics.Elapsed()
			ge.resultCh <- model.GameResult{
				WinnerID:       resp.UserID,
				TotalCorrect:   c,
				TotalIncorrect: ic,
				TimeTaken:      elapsed,
			}
			fmt.Println()
			fmt.Println("========================================")
			fmt.Printf("  WINNER: %s\n", resp.UserID)
			fmt.Printf("  Correct answers:   %d\n", c)
			fmt.Printf("  Incorrect answers: %d\n", ic)
			fmt.Printf("  Time to winner:    %v\n", elapsed)
			fmt.Println("========================================")
			fmt.Println()
			ge.cancel()
		})
	} else {
		ge.metrics.RecordIncorrect()
	}
}

func (ge *gameEngine) drainRemaining() {
	for {
		select {
		case resp, ok := <-ge.responses:
			if !ok {
				return
			}
			if resp.IsCorrect {
				ge.metrics.RecordCorrect()
			} else {
				ge.metrics.RecordIncorrect()
			}
		default:
			return
		}
	}
}
