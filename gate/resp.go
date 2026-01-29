package gate

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/saravanan611/base/log"
)

const (
	Success = "S"
)

type RespStruct struct {
	Status   string `json:"status,omitempty"`
	ErrCode  string `json:"code,omitempty"`
	Msg      string `json:"msg,omitempty"`
	RespInfo any    `json:"info,omitempty"`
}

/* standard error responce structure for api */

func ErrorSender(w http.ResponseWriter, pErrCode string, pErr error) {
	log.Err(pErr)
	log.Info("ErrorSender (+)")
	w.WriteHeader(http.StatusInternalServerError)
	if _, lErr := fmt.Fprintf(w, "Error: << %s >>. Please refer to this code for developer fast support: (%s).", pErr.Error(), pErrCode); lErr != nil {
		log.Err(lErr)
	}
	log.Info("ErrorSender (-)")
}

func MsgSender(w http.ResponseWriter, pInfo any) {
	log.Info("MsgSender (+)")
	log.Info("%+v", pInfo)
	if lErr := json.NewEncoder(w).Encode(RespStruct{Status: Success, RespInfo: pInfo}); lErr != nil {
		log.Err(lErr)
	}
	log.Info("MsgSender (-)")
}
