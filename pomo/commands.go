package pomo

import (
	"fmt"
	"time"

	logging "github.com/azema-shaik/logger"
)

var cmdLogger *logging.Logger

func Start(command *startCmd) *state {

	timeFormat := "02-01-2006 03:04:05.000 PM"

	duration := map[string]time.Duration{
		"hour":   time.Hour,
		"minute": time.Minute,
	}[command.unit]

	timeNow := time.Now()
	userTimeLimit := timeNow.Add(time.Duration(command.duration) * duration)
	stateTimer := NewState(time.Duration(command.duration)*duration,
		time.Duration(command.reminderTime)*time.Minute)

	cmdLogger.Info(fmt.Sprintf("Start At: %v, End At: %v, stateTimers initialized.",
		timeNow, userTimeLimit))

	fmt.Printf("\033[38;5;14mStarting At: %s, Ending At: %s, reminder in: %d minutes, breakTime: %d minutes\033[0m\n",
		timeNow.Format(timeFormat), userTimeLimit.Format(timeFormat), command.reminderTime, command.breakTime)

pomo:
	for {
		select {
		case <-stateTimer.timer.C:
			stateTimer.Stop()
			fmt.Printf("\033[38;5;10mYou have completed %d %s(s) sucessfully\033[0m\n",
				command.duration, command.unit)
			fmt.Printf("\033[38;5;14mTake a break for %d minutes.\033[0m\n", command.breakTime)

			<-time.After(time.Duration(command.breakTime) * time.Minute)
			fmt.Printf("\033[38;5;10mYou have succesfully completed a pomodoro session🥳\033[0m\n")
			break pomo

		case ticker := <-stateTimer.ticker.C:
			timeLeft := userTimeLimit.Sub(ticker)
			if timeLeft.Minutes() < float64(command.reminderTime) {
				fmt.Printf("\033[38;5;14mYou have less than 2 minutes left: %.f\033[0m\n", timeLeft.Minutes())
			} else {
				fmt.Printf("\033[38;5;10mReminder Time: [%s],TimeLeft: [%.0f]\033[0m\n",
					ticker.Format(timeFormat), timeLeft.Minutes())
			}
		}
	}

	return stateTimer

}

func init() {
	cmdLogger = NewLogger("commands")
}
