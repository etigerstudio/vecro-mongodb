package main

import "log"

// Main stressing entry
func stress(t ServiceType) {
	if t == vanilla {
		log.Println("vanilla-stress")
		return
	} else if t == cpu {
		log.Println("cpu-stress")
	} else if t == io {
		log.Println("io-stress")
	}
}