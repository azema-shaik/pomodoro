package main

import (
	pomo "github.com/azema-shaik/pomo/pomo"
)

func main() {

	db := pomo.GetDB()

	cmdType, params := pomo.Cli()
	switch cmdType {
	case pomo.START:
		_ = pomo.Start(db, params)
	case pomo.STATUS:
		pomo.Status(db, params)
	case pomo.RESET:
		pomo.Reset(db, params)

	}

}
