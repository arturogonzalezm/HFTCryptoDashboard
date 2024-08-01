package util

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Logging Helpers

// LogInfo logs informational messages
func LogInfo(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}

// LogError logs error messages
func LogError(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format, args...)
}

// LogAndExit logs the provided message and exits the application with the specified status code
func LogAndExit(format string, code int, args ...interface{}) {
	log.Printf("[FATAL] "+format, args...)
	os.Exit(code)
}

// Error Handling

// CheckError logs the error if not nil
func CheckError(err error, format string, args ...interface{}) {
	if err != nil {
		LogError(format, args...)
	}
}

// Configuration Helpers

// FileExists checks if a file exists
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// ValidateConfig checks if the required config fields are set
func ValidateConfig(cfg interface{}) error {
	// Implement specific validation logic based on the config structure
	// This is a placeholder for a detailed validation logic
	return nil
}

// String Manipulation

// ToUpperCase converts a string to upper case
func ToUpperCase(s string) string {
	return strings.ToUpper(s)
}

// ToLowerCase converts a string to lower case
func ToLowerCase(s string) string {
	return strings.ToLower(s)
}

// Contains checks if a substring is present in a string
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// HTTP Helpers

// MakeGetRequest makes an HTTP GET request and returns the response body
func MakeGetRequest(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			LogError("Failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP request failed with status %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// Time Utilities

// CurrentTimestamp returns the current timestamp in a specified format
func CurrentTimestamp() string {
	return time.Now().Format(time.RFC3339)
}

// ParseTimestamp parses a timestamp string into a time.Time object
func ParseTimestamp(timestamp string, format string) (time.Time, error) {
	return time.Parse(format, timestamp)
}

// Retry attempts to execute the provided function fn up to maxAttempts times,
// with a specified delay between each attempt
func Retry(maxAttempts int, delay time.Duration, fn func() error) error {
	var err error
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		err = fn()
		if err == nil {
			return nil
		}
		LogError("Attempt %d failed: %v. Retrying in %v...", attempts, err, delay)
		time.Sleep(delay)
	}
	return err
}
