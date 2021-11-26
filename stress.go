package main

import (
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"
)

const (
	stressNgName = "./stress-ng"
	cpuOpsBase = 6
	ioOpsBase = 80
)

// Main stressing entry
func stress(delayDuration, delayJitter, cpuLoad, ioLoad int) {
	if delayDuration > 0 {
		delay(delayDuration, delayJitter)
	}

	// TODO: Implement all-in-one stressing rather than individual
	if cpuLoad > 0 {
		cpuStress(cpuLoad)
	}

	if ioLoad > 0 {
		ioStress(ioLoad)
	}

	// TODO: Implement memory stress
}

func delay(delayDuration int, delayJitter int) {
	jitter := math.Floor(rand.Float64() * 2 - 1) * float64(delayJitter)
	time.Sleep(time.Millisecond * time.Duration(delayDuration + int(jitter)))
	log.Printf("slept for %d miliseconds\n", delayDuration + int(jitter))
}

func cpuStress(cpuLoad int) {
	// Approximately 6,000 ops per sec => 6 ops per 1ms
	// on Intel(R) Xeon(R) CPU E5-2630 v4 @ 2.20GHz
	execCommand(stressNgName,
		"--cpu", "1",
		"--cpu-ops", strconv.Itoa(cpuOpsBase * cpuLoad),
		"--cpu-method", "sqrt")
	log.Printf("cpu load amount: %d\n", cpuLoad)
}

func ioStress(ioLoad int) {
	// Approximately 80,000 ops per sec => 80 ops per 1ms
	// on Intel(R) Xeon(R) CPU E5-2630 v4 @ 2.20GHz
	execCommand(stressNgName,
		"--io", "1",
		"--io-ops", strconv.Itoa(ioOpsBase * ioLoad))
	log.Printf("io load amount: %d\n", ioLoad)
}