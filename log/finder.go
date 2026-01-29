package log

import (
	"log"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var (
	gRequestIDMap = sync.Map{}
)

// Get current goroutine ID
func getGoroutineID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = b[len("goroutine "):]
	b = b[:strings.IndexByte(string(b), ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

// Set request ID for current goroutine
func SetRequestID(requestID string) {
	gid := getGoroutineID()
	gRequestIDMap.Store(gid, requestID)
}

// Get request ID for current goroutine
func GetRequestID() string {
	gid := getGoroutineID()
	log.Println("gid :", gid)
	if id, ok := gRequestIDMap.Load(gid); ok {
		return id.(string)
	}
	return "UNKNOWN"
}

// Clear request ID when done
func ClearRequestID() {
	gid := getGoroutineID()
	gRequestIDMap.Delete(gid)
}
