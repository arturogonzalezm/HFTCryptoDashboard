package util

import (
	"log"
	"os"
	"time"
)

// Retry attempts to execute the provided function fn up to maxAttempts times,
// with a specified delay between each attempt.
func Retry(maxAttempts int, delay time.Duration, fn func() error) error {
	var err error
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		err = fn()
		if err == nil {
			return nil
		}
		log.Printf("Attempt %d failed: %v. Retrying in %v...", attempts, err, delay)
		time.Sleep(delay)
	}
	return err
}

// LogAndExit logs the provided message and exits the application with the specified status code.
func LogAndExit(format string, code int, args ...interface{}) {
	log.Printf(format, args...)
	os.Exit(code)
}

// CurrentTimestamp returns the current timestamp in a specified format.
func CurrentTimestamp() string {
	return time.Now().Format(time.RFC3339)
}
