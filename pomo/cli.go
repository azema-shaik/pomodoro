package pomo

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"

	logging "github.com/azema-shaik/logger"
)

type cmd int

const (
	_ = iota
	START
	STATUS
	RESET
)

const validMainCommands = "start|status"

var cliLogger *logging.Logger

func validate(arg string) (valid bool, message string) {
	validCommands := strings.Split(validMainCommands, "|")
	valid = slices.Contains(validCommands, arg)
	if !valid {
		message = fmt.Sprintf("Invalid main command.[MAIN_COMMAND] should be one of {%s}", validMainCommands)
	}

	return

}

func checkIfHelpExists(args []string, cmd *flag.FlagSet) {
	for _, helpFlag := range []string{"-help", "-h", "--help"} {
		if slices.Contains(args, helpFlag) {
			cmd.PrintDefaults()
			os.Exit(0)
		}
	}
}

func check(duration, checkTime int, time_unit string) bool {
	return checkTime < 0 && time_unit == "minute" && duration < checkTime
}

func startCmd(args []string) (params map[string]any) {
	var duration, breakTime, reminderTime int
	var unit, session_name string

	start := flag.NewFlagSet("start", flag.ExitOnError)
	start.IntVar(&duration, "duration", 25, "Number of times to concentrate. Defaults to 25.")
	start.StringVar(&unit, "unit", "minute", "Duration unit can only be one of \"minute\" or \"hour\".Defaults to minute")
	start.IntVar(&breakTime, "break", 5, "Break time. Defaults to 5 minutes")
	start.IntVar(&reminderTime, "reminder", 2, "Remind time.Defaults to 2 minutes.")
	start.StringVar(&session_name, "name", "", "this is session name and is always required.")

	start.Parse(args)
	checkIfHelpExists(args, start)
	if session_name == "" {
		fmt.Println("\033[38;5;9m[INFO]:[session name cannot be empty]\033[0m")
		os.Exit(1)
	}

	cliLogger.Debug(fmt.Sprintf("Duration = %d,Unit = %s,breakTime = %d", duration, unit, breakTime))

	if duration < 1 {
		message := fmt.Sprintf("Invalid argument for \"duration\": %d", duration)
		cliLogger.Error(message)
		fmt.Println("\033[38;5;9m[INFO]: [Invalid value for duration initalizing it to 25]\033[0m\n")

	}

	if !slices.Contains([]string{"minute", "hour"}, strings.TrimRight(unit, "s")) {
		cliLogger.Error(fmt.Sprintf("Invalid argument for \"unit\": \"%s\"", unit))
		fmt.Printf("\033[38;5;9m[INFO]: [invalid value for time unit initializing it to \"minute\"]")

	}

	if check(duration, reminderTime, unit) {
		message := fmt.Sprintf("Invalid value for \"reminder\": %d", reminderTime)
		cliLogger.Error(message)
		fmt.Println("[INFO]: [reminder will be iniitalized to 5]")
		reminderTime = 5

	}

	if check(duration, breakTime, unit) {
		message := fmt.Sprintf("Invalid value for \"break\": %d", breakTime)
		cliLogger.Error(message)
		fmt.Println("[INFO]: [breaktime will be iniitalized to 5]")
		breakTime = 5
	}

	cliLogger.Info(fmt.Sprintf("Start main command initalized with duration = %d and unit=%s and break = %d", duration, unit, breakTime))
	return map[string]any{"session_name": session_name,
		"duration":   duration,
		"break_time": breakTime,
		"reminder":   reminderTime,
		"time_unit":  unit}
}

func statusAndReset(args []string, cmdName string) (params map[string]any) {
	var name string

	cmd := flag.NewFlagSet(cmdName, flag.ExitOnError)
	cmd.StringVar(&name, "name", "all", fmt.Sprintf("fetches %s by name", cmdName))

	cmd.Parse(args)
	checkIfHelpExists(args, cmd)
	return map[string]any{"session_name": name}

}

func init() {
	cliLogger = NewLogger("cli")
}

func Cli() (cmdType cmd, params map[string]any) {
	if len(os.Args) < 2 {
		fmt.Printf("\033[38;5;9mNo valid subcommand found.Should be one of %s\033[0m\n", validMainCommands)
		os.Exit(1)
	}

	mainCommand := strings.ToLower(os.Args[1])
	if isValid, message := validate(mainCommand); !isValid {
		fmt.Printf("\033[38;5;9m%s\033[0m\n", message)
		os.Exit(1)
	}

	args := map[bool][]string{true: []string{}, false: os.Args[2:]}[len(os.Args) < 3]
	cliLogger.Debug(fmt.Sprintf("cli arguments: %v", args))
	switch mainCommand {
	case "start":
		params = startCmd(args)
		cmdType = START
	case "status":
		params = statusAndReset(args, "status")
		cmdType = STATUS
	case "reset":
		params = statusAndReset(args, "reset")
		cmdType = RESET

	}

	cliLogger.Info(fmt.Sprintf("command select: %s, params: %v",
		map[cmd]string{START: "start", STATUS: "status", RESET: "reset"}[cmdType], params))
	return cmdType, params

}
