package common

import (
	"log"
	"runtime"
	"strconv"
)

func SetMaxProcs() {
	maxProcsRaw := GetMaxProcs()

	if maxProcsRaw != "" {
		maxProcs, err := strconv.Atoi(maxProcsRaw)
		if err != nil {
			log.Fatal(err.Error())
		}

		runtime.GOMAXPROCS(maxProcs)
	}
}
