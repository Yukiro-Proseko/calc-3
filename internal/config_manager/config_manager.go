package config_manager

import (
	"os"
	"strconv"
)

var (
	ComputingPwr, _     = strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	AdditionMS, _       = strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
	SubtractionMS, _    = strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS"))
	MultiplicationMS, _ = strconv.Atoi(os.Getenv("TIME_MULTIPLICATION_MS"))
	DivisionMS, _       = strconv.Atoi(os.Getenv("TIME_DIVISION_MS"))
)
