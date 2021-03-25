package utils

import (
	"fmt"
	"time"
)

// TimeTrack starts a deferred timer, in order to profile the execution of a function
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("Scraping %s took %s\n", name, elapsed)
}
