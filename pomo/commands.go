package pomo

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	logging "github.com/azema-shaik/logger"
)

var cmdLogger *logging.Logger

func Start(params map[string]any) *state {

	timeUnit := params["time_unit"].(string)
	duration := params["duration"].(int)
	reminderTime := params["reminder"].(int)
	breakTime := params["break_time"].(int)

	timeFormat := "02-01-2006 03:04:05.000 PM"

	timeDuration := map[string]time.Duration{
		"hour":   time.Hour,
		"minute": time.Minute,
	}[timeUnit]

	timeNow := time.Now()
	userTimeLimit := timeNow.Add(time.Duration(duration) * timeDuration)
	stateTimer := NewState(time.Duration(duration)*timeDuration,
		time.Duration(reminderTime)*time.Minute)

	cmdLogger.Info(fmt.Sprintf("Start At: %v, End At: %v, stateTimers initialized.",
		timeNow, userTimeLimit))

	fmt.Printf("\033[38;5;14mStarting At: %s, Ending At: %s, reminder in: %d minutes, breakTime: %d minutes\033[0m\n",
		timeNow.Format(timeFormat), userTimeLimit.Format(timeFormat),
		reminderTime, breakTime)

	for {
		select {
		case <-stateTimer.timer.C:
			stateTimer.Stop()
			fmt.Printf("\033[38;5;10mYou have completed %d %s(s) sucessfully\033[0m\n",
				duration, timeUnit)
			fmt.Printf("\033[38;5;14mTake a break for %d minutes.\033[0m\n", breakTime)

			<-time.After(time.Duration(breakTime) * time.Minute)
			fmt.Printf("\033[38;5;10mYou have succesfully completed a pomodoro session🥳\033[0m\n")
			return stateTimer

		case ticker := <-stateTimer.ticker.C:
			timeLeft := userTimeLimit.Sub(ticker)
			if timeLeft.Minutes() < float64(reminderTime) {
				fmt.Printf("\033[38;5;14mYou have less than 2 minutes left: %.f\033[0m\n", timeLeft.Minutes())
			} else {
				fmt.Printf("\033[38;5;10mReminder Time: [%s],TimeLeft: [%.0f]\033[0m\n",
					ticker.Format(timeFormat), timeLeft.Minutes())
			}
		}
	}

}

func Status(db *sql.DB, params map[string]any) {
	session_name := params["session_name"].(string)
	query := map[bool]string{true: `SELECT session_name, timestamp, time_duration, break_time, time_unit,reminder_time FROM pomo
									WHERE username = :username`,
		false: `SELECT session_name, timestamp, time_duration, break_time, time_unit,reminder_time FROM pomo
									WHERE username = :username AND session_name = :session_name`}[session_name == "all"]
	cmdLogger.Debug(fmt.Sprintf("status sql command: %s", query))

	countRows := 0

	rows := fetchRows(db, query, session_name)

	var totalTime int
	for rows.Next() {
		var unit, session_name string
		var session_duration, session_break, session_reminder int
		var timestamp time.Time

		err := rows.Scan(&session_name, &timestamp, &session_duration, &session_break, &unit, session_reminder)
		if err != nil {
			cmdLogger.Error(fmt.Sprintf("error when fetching rows. err: %s", err.Error()))
			fmt.Println("\033[38;5;9mIssue when fetching errors.Consult logs.\033[0m")
			os.Exit(1)
		}

		fmt.Printf("\033[38;5;10mSession Name: %s, Session Time: %s\nSession Duration: %d %s, Session Break: %d minute(s), Session Reminder: %d minute(s)\033[0m",
			session_name, utcToIst(timestamp), session_duration, unit, session_break, session_reminder)
		fmt.Println(strings.Repeat("=", 50))

		totalTime += map[string]int{
			"hour":   session_duration * 60,
			"minute": session_duration}[unit]

		countRows += 1

	}

	if countRows == 0 {
		fmt.Println("\033[38;5;14mNo records found.\033[0m")
	} else {
		duration := map[bool]time.Duration{
			true:  time.Hour,
			false: time.Minute}[totalTime > 60]
		fmt.Printf("\033[38;5;14mBe proud you have completed %s\033[0m\n", time.Duration(totalTime)*duration)
	}
}

func Reset(db *sql.DB, params map[string]any) {
	session_name := params["session_name"].(string)
	query := map[bool]string{
		true:  `DELETE FROM pomo WHERE username = :username RETURNING session_name, timestamp, time_duration, time_unit`,
		false: `DELETE FROM pomo WHERE username = :username AND session_name = :name RETURNING session_name, timestamp, time_duration, time_unit`,
	}[session_name == "all"]

	rows := fetchRows(db, query, session_name)
	w := tabwriter.NewWriter(os.Stdout, 5, 5, 5, ' ', tabwriter.AlignRight|tabwriter.Debug)
	defer w.Flush()

	fmt.Fprintln(w, "session_name\ttimestamp\ttime_duration")
	rowCount := 0
	for rows.Next() {
		var session_name, time_unit string
		var time_duration int
		var timestamp time.Time

		err := rows.Scan(&session_name, &timestamp, &time_duration, &time_unit)
		if err != nil {
			cmdLogger.Error(fmt.Sprintf("error when scaning rows: %s", err.Error()))
			fmt.Println("\033[38;5;9mError when trying initalizing app. Consult logs.\033[0m")
			os.Exit(1)
		}

		fmt.Fprintf(w, "%s\n%s\t%d %s\n", session_name, utcToIst(timestamp), time_duration, time_unit)
		rowCount += 1
	}

	if rowCount == 0 {
		fmt.Println("\033[38;5;14mNo rows found.\033[0m")
	}

}

func init() {
	cmdLogger = NewLogger("commands")
}
