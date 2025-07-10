package pomo

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	logging "github.com/azema-shaik/logger"
	_ "modernc.org/sqlite"
)

var dataLogger *logging.Logger

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

func createDB(db *sql.DB) error {
	dataLogger.Info("Startint to create table.")
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS pomo(
	username TEXT,
	session_name TEXT, 
	time_duration INTEGER,
	break_time INTEGER,
	reminder_time INTEGER,
	time_unit TEXT,
	timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	return err
}

func GetDB() *sql.DB {

	config, isExists := checkConfigExists()
	cmdLogger.Info(fmt.Sprintf("config: %v, isExists: %v", config, isExists))
	dbName := map[bool]string{
		false: "pomo.db",
		true:  config["db"]}[isExists]

	cmdLogger.Info(fmt.Sprintf("DbName: %s", dbName))

	db, err := sql.Open("sqlite", filepath.Join(path, dbName))
	cmdLogger.Debug(fmt.Sprintf("db after initalization: %#v", db))
	if err != nil {
		dataLogger.Error(fmt.Sprintf("Error when trying to connect database. err: %s", err.Error()))
		fmt.Println("error when trying to connecting to db. Please consult logs.")
		os.Exit(1)
	}

	if _, err = db.Query("SELECT * FROM pomo"); err != nil && strings.Contains(err.Error(), "no such table") {
		dataLogger.Info(fmt.Sprintf("table does not exist. Creating table. Err: %s", err.Error()))
		if err := createDB(db); err != nil {
			dataLogger.Error(fmt.Sprintf("error when trying to create table. err: %s", err.Error()))
			fmt.Println("\033[38;5;9mError when trying to initalize app.Consult logs.\033[0m")
			os.Exit(1)
		}
		return db
	}

	return db
}

func update(db *sql.DB, params map[string]any) {
	stmt, err := db.Prepare(`INSERT INTO pomo (username,session_name, time_duration,
	break_time, reminder_time, time_unit) VALUES 
	(:username, :session_name, :time_duration, :break_time,:reminder_time,
		:time_unit)`)
	if err != nil {
		dataLogger.Error(fmt.Sprintf("Error when trying to update db. Err: %s", err.Error()))
		//try thinkin of something concurrent to save state temporarily and then when the app starts again sweep the local saved.
		fmt.Println("\033[38;5;9merror when trying to update table. please consult logs.\033[0m")
		os.Exit(1)
	}

	_, err = stmt.Query(sql.Named("username", os.Getenv("USERNAME")),
		sql.Named("session_name", params["session_name"]),
		sql.Named("time_duration", params["duration"]),
		sql.Named("break_time", params["break_time"]),
		sql.Named("reminder_time", params["reminder"]),
		sql.Named("time_unit", params["time_unit"]))
	if err != nil {
		dataLogger.Error(fmt.Sprintf("Error when trying to update db. Err: %s", err.Error()))
		fmt.Println("\033[38;5;9merror when trying to update table. please consult logs.\033[0m")
		os.Exit(1)
	}

}

func init() {
	dataLogger = NewLogger("data")
}
