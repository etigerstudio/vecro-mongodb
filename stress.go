package main

import (
	"log"
	"math"
	"math/rand"
	"time"
)

)

// Main stressing entry
func stress(delayTime, delayJitter, cpuLoad, ioLoad int) {
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