package pomodoro

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

func start(args []string) (duration int, unit string) {
	start := flag.NewFlagSet("start", flag.ExitOnError)
	start.IntVar(&duration, "duration", 25, "Number of times to concentrate. Defaults to 25.")
	start.StringVar(&unit, "unit", "minute", "Duration unit can only be one of \"minute\" or \"hour\".Defaults to minute")

	start.Parse(args)

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
	cliLogger.Info(fmt.Sprintf("Start main command initalized with duration = %d and unit=%s", duration, unit))
	return duration, unit

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

	switch mainCommand {
	case "start":
		start(args[2:])
	case "status":
		//to be implemeted.
	}

}
