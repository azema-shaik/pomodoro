package pomo

import "time"

type state struct {
	timer  *time.Timer
	ticker *time.Ticker
}

func NewState(timerDuration time.Duration,
	tickerDuration time.Duration) *state {

	return &state{
		timer:  time.NewTimer(timerDuration),
		ticker: time.NewTicker(tickerDuration),
	}

}

func (s *state) Stop() {
	s.ticker.Stop()
	s.timer.Stop()
}
