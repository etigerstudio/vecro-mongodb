package main

import (
	"log"
	"math"
	"math/rand"
	"time"
)

const (
	delayTimeEnvKey = "BEN_DELAY_TIME"
	delayJitterEnvKey = "BEN_DELAY_JITTER"
	cpuLoadEnvKey = "BEN_CPU_WORKLOAD"
	ioLoadEnvKey = "BEN_IO_WORKLOAD"
)

var (
	delayTime int
	delayJitter int
	cpuLoad int
	ioLoad int
)

func init() {
	delayTime, _ = getEnvInt(delayTimeEnvKey, 0)
	delayJitter, _ = getEnvInt(delayJitterEnvKey, delayTime / 10)
	cpuLoad, _ = getEnvInt(cpuLoadEnvKey, 0)
	ioLoad, _ = getEnvInt(ioLoadEnvKey, 0)
}

// Main stressing entry
func stress() {
	if delayTime > 0 {
		delay(delayTime, delayJitter)
	}

	if cpuLoad > 0 {
		cpuStress(cpuLoad)
	}

	if ioLoad > 0 {
		ioStress(ioLoad)
	}

	// TODO: Implement memory stress
}

func delay(delayTime int, delayJitter int) {
	jitter := math.Floor(rand.Float64() * 2 - 1) * float64(delayJitter)
	time.Sleep(time.Millisecond * time.Duration(delayTime + int(jitter)))
	log.Printf("slept for %d miliseconds\n", delayTime+int(jitter))
}

func cpuStress(cpuLoad int) {
	log.Printf("cpu load amount: %d\n", cpuLoad)
}

func ioStress(ioLoad int) {
	log.Printf("io load amount: %d\n", ioLoad)
}