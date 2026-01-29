package gate

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/saravanan611/base/log"
)

/*
=========================================================================================================================
restart the program on every day mention time to cleare the in-memory/catch-memory for good pratice and free the resource
=========================================================================================================================
*/
func AutoRestart(lRestartHour, lRestartMinute int) {
	log.Info("AutoRestart (+)")

	go func() {

		lCurrentTime := time.Now()
		if !(lCurrentTime.Hour() == lRestartHour && lCurrentTime.Minute() == lRestartMinute) {
			lNextExecutionTime := time.Date(lCurrentTime.Year(), lCurrentTime.Month(), lCurrentTime.Day(), lRestartHour, lRestartMinute, 0, 0, lCurrentTime.Location())
			if lNextExecutionTime.Before(lCurrentTime) {
				lNextExecutionTime = lNextExecutionTime.Add(time.Duration(24 * time.Hour))
			}
			fmt.Println("Current Time:", lCurrentTime)
			fmt.Println("Next Execution Time:", lNextExecutionTime, lNextExecutionTime.Sub(lCurrentTime))
			durationUntilNextExecution := lNextExecutionTime.Sub(lCurrentTime)
			time.Sleep(durationUntilNextExecution)
		}

		log.Info("program going to restart within a minute...")
		time.Sleep(1 * time.Minute)

		lErr := restart()

		if lErr != nil {
			log.Err(lErr)
		}
		os.Exit(0)
		log.Info("AutoRestart (-)")
	}()
}

func restart() (lErr error) {
	log.Info("restart (+)")
	execPath, lErr := os.Executable()
	if lErr != nil {
		return lErr
	}

	cmd := exec.Command(execPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	lErr = cmd.Start()
	if lErr != nil {
		return lErr
	}
	log.Info("restart (-)")
	return
}

/*
===================================================================================================================
this function treager in end of the program  like close glogal db connection
use defer to call this in func main() , this func is also auto restart your code where your program will panic
===================================================================================================================
*/

func TreagerOnEnd(pEndFunc ...func()) {
	log.Info("TreagerOnEnd (+)")

	for lIdx, lFunc := range pEndFunc {
		log.Debug("going to execuate the end process %d \n", lIdx)
		lFunc()
	}

	if lErr := recover(); lErr != nil {
		log.Err(lErr)
		if lErr := restart(); lErr != nil {
			log.Err(lErr)
		}
	}

	log.Info("TreagerOnEnd (-)")
}
