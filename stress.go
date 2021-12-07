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
	cpuOpsBase = 3
	ioOpsBase = 4
)

// Main stressing entry
func stress(delayTime, delayJitter, cpuLoad, ioLoad int) {
	if delayTime > 0 {
		delay(delayTime, delayJitter)
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

func delay(delayTime int, delayJitter int) {
	jitter := math.Floor((rand.Float64() * 2 - 1) * float64(delayJitter))
	time.Sleep(time.Millisecond * time.Duration(delayTime + int(jitter)))
	log.Printf("slept for %d miliseconds\n", delayTime + int(jitter))
}

func cpuStress(cpuLoad int) {
	// Approximately 3,000 ops per sec => 3 ops per 1ms
	// on Intel(R) Xeon(R) CPU E5-2630 v4 @ 2.20GHz
	execCommand(stressNgName,
		"--cpu", "1",
		"--cpu-ops", strconv.Itoa(cpuOpsBase * cpuLoad),
		"--cpu-method", "sqrt")
	log.Printf("cpu load amount: %d\n", cpuLoad)
}

func ioStress(ioLoad int) {
	// Approximately 4,000 ops per sec => 4 ops per 1ms
	// on Intel(R) Xeon(R) CPU E5-2630 v4 @ 2.20GHz
	execCommand(stressNgName,
		"--hdd", "1",
		"--hdd-ops", strconv.Itoa(ioOpsBase * ioLoad),
	)
	log.Printf("io load amount: %d\n", ioLoad)
}