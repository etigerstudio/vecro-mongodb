package main

import (
	"bufio"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	cpuOpsBase = 16384 * 2 // stress-ng use 16384 as a minimal sqrt(rand()) unit
	ioOpsBase = 64 // 64 times * 1 KiB = 64 KiB
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
	// Approximately 16,384 * 2 ops per 1ms
	// on Intel(R) Xeon(R) CPU E5-2630 v4 @ 2.20GHz
	for i := 0; i < cpuOpsBase * cpuLoad; i++ {
		math.Sqrt(rand.Float64())
	}
	log.Printf("cpu load amount: %d, total sqrt(rand()): %d\n", cpuLoad, cpuOpsBase * cpuLoad)
}

// FIXME: Better file specific error handling in this function
func ioStress(ioLoad int) {
	filename := "/tmp/ben_base_stress_" + strconv.FormatUint(rand.Uint64(), 10)
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriter(f)
	totalBytes, currentBytes := 0, 0
	for i := 0; i < ioOpsBase * ioLoad; i++ {
		// Write 1 KiB to buffer at one time
		n, err := w.WriteString(strings.Repeat("contents", 128))
		if err != nil {
			panic(err)
		}

		totalBytes += n
		currentBytes += n
		// Flush buffered data every 512 KiB
		if currentBytes >= 1024 * 512 {
			err = w.Flush()
			if err != nil {
				panic(err)
			}
			err = f.Sync()
			if err != nil {
				panic(err)
			}

			currentBytes -= 1024 * 512
		}
	}

	err = w.Flush()
	if err != nil {
		panic(err)
	}
	err = f.Sync()
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}

	err = os.Remove(filename)
	if err != nil {
		panic(err)
	}
	log.Printf("io load amount: %d, total bytes written: %d\n", ioLoad, ioOpsBase * 8 * 128)
}