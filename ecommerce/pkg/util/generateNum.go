package util

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

var (
    numberMutex sync.Mutex
    lastTime    int64
)

// GenerateNo generates a unique number with given prefix
// Format: PREFIX + YYYYMMDDHHMMSS + 6-digit random number
func GenerateNo(prefix string) string {
    numberMutex.Lock()
    defer numberMutex.Unlock()

    now := time.Now()
    timestamp := now.Format("20060102150405")
    currentTime := now.UnixNano() / 1e6

    // Ensure unique timestamp
    if currentTime <= lastTime {
        currentTime = lastTime + 1
    }
    lastTime = currentTime

    // Generate 6-digit random number
    random := rand.New(rand.NewSource(currentTime))
    randomNum := random.Intn(1000000)

    return fmt.Sprintf("%s%s%06d", prefix, timestamp, randomNum)
}