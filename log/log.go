package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

// ==================================================
// Global map to store requestID for each goroutine
// ==================================================
var (
	GEnvVal bool
	gInfo   = log.New(os.Stdout, "INFO: ", log.LstdFlags|log.Lshortfile)
	gDebug  = log.New(os.Stdout, "DEBUG: ", log.LstdFlags|log.Lshortfile)
	gErr    = log.New(os.Stderr, "ERROR: ", log.LstdFlags|log.Lshortfile)
	gEnErr  = log.New(os.Stderr, "ERROR: ", log.LstdFlags)
)

/* -----------------------custom error ------------------- */
type ownErr struct {
	lFileInfo string
	lErr      string
}

/* -----------------------err to string-------------------- */
func (pErr *ownErr) Error() string {
	return pErr.lErr
}

// ---------- ERROR WRAPPER ----------
func Error(pErr any) error {

	if lErr, lOk := pErr.(*ownErr); lOk {
		return lErr
	}

	_, lFile, lLine, _ := runtime.Caller(1)
	lStrArray := strings.Split(lFile, "/")
	lFilename := lStrArray[len(lStrArray)-2] + "/" + lStrArray[len(lStrArray)-1]
	return &ownErr{lFileInfo: fmt.Sprintf("%s:%d", lFilename, lLine), lErr: fmt.Sprintf("%v", pErr)}

}

// ============================================
// Logging functions - NO PARAMETERS NEEDED!
// ============================================

// ---------- INFO LOGGER ----------
func Info(format string, args ...any) {
	if GEnvVal {
		return
	}
	requestID := GetRequestID()
	msg := fmt.Sprintf(format, args...)
	gInfo.Output(2, fmt.Sprintf("[ReqID: %s] %s", requestID, msg))
}

// ---------- Debug LOGGER ----------
func Debug(format string, args ...any) {
	if GEnvVal {
		return
	}
	requestID := GetRequestID()
	msg := fmt.Sprintf(format, args...)
	gDebug.Output(2, fmt.Sprintf("[ReqID: %s] %s", requestID, msg))
}

// ---------- Err LOGGER ----------
func Err(pErr any) {
	requestID := GetRequestID()
	if lErr, lok := pErr.(*ownErr); lok {
		gEnErr.Printf("%s [ReqID: %s] %s", lErr.lFileInfo, requestID, lErr.Error())
		return
	}
	// use Output(counter+2, ...) so call depth aligns properly
	gErr.Output(2, fmt.Sprintf("[ReqID: %s] %v", requestID, pErr))
}
