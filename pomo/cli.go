package pomo

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"

	logging "github.com/azema-shaik/logger"
)

const validMainCommands = "start|status"

var cliLogger *logging.Logger

func validate(arg string) (valid bool, message string) {
	validCommands := strings.Split(validMainCommands, "|")
	valid = slices.Contains(validCommands, arg)
	if valid {
		message = fmt.Sprintf("Invalid main command.[MAIN_COMMAND] should be one of {%s}", validMainCommands)
	}

	return

}

func start(args []string) *startCmd {
	var duration, breakTime, reminderTime int
	var unit string

	start := flag.NewFlagSet("start", flag.ExitOnError)
	start.IntVar(&duration, "duration", 25, "Number of times to concentrate. Defaults to 25.")
	start.StringVar(&unit, "unit", "minute", "Duration unit can only be one of \"minute\" or \"hour\".Defaults to minute")
	start.IntVar(&breakTime, "break", 5, "Break time. Defaults to 5 minutes")
	start.IntVar(&reminderTime, "reminder", 2, "Remind time.Defaults to 2 minutes.")

	start.Parse(args)
	cliLogger.Debug(fmt.Sprintf("Duration = %d,Unit = %s,breakTime = %d", duration, unit, breakTime))

	if duration < 1 {
		message := fmt.Sprintf("Invalid argument for \"duration\": %d", duration)
		cliLogger.Error(message)
		fmt.Printf("\033[38;5;9m%s\033[0m\n", message)
		os.Exit(1)
	}
	if !slices.Contains([]string{"minute", "hour"}, strings.TrimRight(unit, "s")) {
		cliLogger.Error(fmt.Sprintf("Invalid argument for \"unit\": \"%s\"", unit))
		fmt.Println("\033[38;5;9mInvalid argument for hour.\033[0m")
		os.Exit(1)
	}
	if breakTime < 0 {
		message := fmt.Sprintf("Invalid argument for \"break\": %d", breakTime)
		cliLogger.Error(message)
		fmt.Printf("\033[38;5;9m%s\033[0m\n", message)
		os.Exit(1)
	}
	if unit == "minute" && (duration < breakTime) {
		message := fmt.Sprintf("break(%d) cannot be greater than duration(%d).", breakTime, duration)
		cliLogger.Error(message)
		fmt.Printf("\033[38;5;9m%s\033[0m\n", message)
		os.Exit(1)
	}

	if reminderTime < 0 {
		message := fmt.Sprintf("Invalid argument for \"reminder\": %d", reminderTime)
		cliLogger.Error(message)
		fmt.Printf("\033[38;5;9m%s\033[0m\n", message)
		os.Exit(1)
	}

	if unit == "minute" && (duration < reminderTime) {
		message := fmt.Sprintf("reminder(%d) cannot be greater than duration(%d).", reminderTime, duration)
		cliLogger.Error(message)
		fmt.Printf("\033[38;5;9m%s\033[0m\n", message)
		os.Exit(1)
	}

	cliLogger.Info(fmt.Sprintf("Start main command initalized with duration = %d and unit=%s and break = %d", duration, unit, breakTime))
	return &startCmd{duration: duration, unit: unit, breakTime: breakTime, reminderTime: reminderTime}
}

func init() {
	cliLogger = NewLogger("cli")
}

func Cli(args []string) {
	if len(args) < 2 {
		fmt.Printf("\033[38;5;9mNo valid subcommand found.Should be one of %s\033[0m\n", validMainCommands)
		os.Exit(1)
	}

	mainCommand := strings.ToLower(args[1])
	if isValid, message := validate(mainCommand); !isValid {
		fmt.Printf("\033[38;5;9m%s\033[0m\n", message)
		os.Exit(1)
	}

	args = args[2:]
	cliLogger.Debug(fmt.Sprintf("cli arguments: %v", args))
	switch mainCommand {
	case "start":
		command := start(args)
		Start(command)
	case "status":
		//to be implemeted.
	}

}
