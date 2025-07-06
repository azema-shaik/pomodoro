package pomo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/azema-shaik/logger"
)

var utilLogger *logger.Logger
var timeZone *time.Location
var timeFormat string = "Monday 02-01-2006 03:04:05 PM MST"
var path string

func utcToIst(timestamp time.Time) string {
	return timestamp.In(timeZone).Format(timeFormat)
}

func fetchRows(db *sql.DB, query, session_name string) *sql.Rows {
	stmt, err := db.Prepare(query)
	if err != nil {
		cmdLogger.Error(fmt.Sprintf("error when trying to prepare query, %s", err.Error()))
		fmt.Println("\033[38;5;9mIssue when trying to initalize connection to the app. Consult logs.\033[0m")
		os.Exit(1)
	}

	rows, err := stmt.Query(sql.Named("username", os.Getenv("USERNAME")),
		sql.Named("session_name", session_name))
	if err != nil {
		cmdLogger.Error(fmt.Sprintf("error when trying to query db. err: %s", err.Error()))
		fmt.Println("\033[38;5;9mError when trying to initialize connection to the app. Consult logs.\033[0m")
		os.Exit(1)
	}
	return rows
}

func checkConfigExists() (config map[string]string, isExist bool) {
	file, err := os.Open(filepath.Join(path, "config.json"))
	if err != nil {
		utilLogger.Info(fmt.Sprintf("config file does not exist"))
		return config, false
	}

	defer file.Close()

	if err := json.NewDecoder(file).Decode(&config); err != nil {
		utilLogger.Error(fmt.Sprintf("Error when tryint to decode json.Check logs.Defaulting.Err: %s", err.Error()))
		return config, false
	}

	return config, true

}

func init() {
	utilLogger = NewLogger("utils")
	timeZone, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		utilLogger.Error(fmt.Sprintf("error initalizing timezone defaulting to utc. err: %s", err))
	}

	timeZone = map[bool]*time.Location{
		true:  time.UTC,
		false: timeZone}[err != nil]

	path, err = os.Getwd()
	if err != nil {
		utilLogger.Error(fmt.Sprintf("error when trying to initliaze filepath. Err: %s\n", err.Error()))
		fmt.Printf("\033[38;5;9mError when trying to inialize filepath. Consult logs\033[0m\n")
		os.Exit(1)
	}

}
