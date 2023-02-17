package common

import (
	"log"
	"os"
)

func GetResourcesPath() string {
	path := os.Getenv("FSTRESSER_RESOURCES_PATH")

	if path == "" {
		log.Fatal("Profiles path not defined")
	}

	return path
}

func GetMaxProcs() string {
	return os.Getenv("FSTRESSER_MAX_PROCS")
}

func GetLogsPath() string {
	path := os.Getenv("FSTRESSER_LOGS_PATH")

	if path == "" {
		log.Fatal("Logs path not defined")
	}

	return path
}
