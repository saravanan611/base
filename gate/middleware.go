package gate

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/saravanan611/base/log"
)

var (
	osSignal        = []os.Signal{}
	allowOrigin     = "*"
	allowCredential = false
	allowHeader     = []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "credentials"}
	methods         = []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace}
)

/*
===============================
set the header for your request
===============================
*/
func SetHeader(pHeader ...string) {
	if len(pHeader) > 0 {
		allowHeader = append(allowHeader, pHeader...)
	}
}

/*
=================================
set signal for greasefull shadown
=================================
*/
func SetSignal(pSignal ...os.Signal) {
	if len(pSignal) > 0 {
		osSignal = append(osSignal, pSignal...)
	}
}

/*
=============================================
set the origin for your request ,default "*""
=============================================
*/
func SetOrigin(pOrigin string) {
	allowOrigin = pOrigin
}

/*
=========================================================================================================
enable Credential to true for cookie_set,some other browser side operation dun by golang ,default "false"
=========================================================================================================
*/
func EnableCredential() {
	allowCredential = true
}

type ResponseCaptureWriter struct {
	http.ResponseWriter
	status int
	body   []byte
}

func (rw *ResponseCaptureWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *ResponseCaptureWriter) Write(body []byte) (int, error) {
	rw.body = append(rw.body, body...)
	return rw.ResponseWriter.Write(body)
}

func (rw *ResponseCaptureWriter) Status() int {
	if rw.status == 0 {
		return http.StatusOK
	}
	return rw.status
}

func (rw *ResponseCaptureWriter) Body() []byte {
	return rw.body
}

// ===============================
// Middleware - Sets up request ID
// ===============================

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(pResp http.ResponseWriter, pReq *http.Request) {
		// Initialize the logger
		(pResp).Header().Set("Access-Control-Allow-Origin", allowOrigin)
		(pResp).Header().Set("Access-Control-Allow-Credentials", fmt.Sprint(allowCredential))
		(pResp).Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
		(pResp).Header().Set("Access-Control-Allow-Headers", strings.Join(allowHeader, ","))

		requestID := strings.ReplaceAll(uuid.New().String(), "-", "")

		log.GEnvVal = strings.EqualFold(os.Getenv("InfoFlog"), "Y")

		// Set request ID for this goroutine
		log.SetRequestID(requestID)
		defer log.ClearRequestID() // Clean up when request completes

		(pResp).Header().Set("X-Request-ID", requestID)
		captureWriter := &ResponseCaptureWriter{ResponseWriter: pResp}
		log.Info("logMiddleware (+)")

		lReqRec := GetRequestorDetail(pReq)
		log.Debug("Method: %s, Path: %s", lReqRec.Method, pReq.URL.Path)
		log.Info("Req Info : %s", lReqRec)

		next.ServeHTTP(captureWriter, pReq)

		log.Debug("Resp Info : %s", string(captureWriter.Body()))

		log.Info("logMiddleware (-)")
	})
}

/*
==============================
set up your server to execuate
==============================
*/
func SetServer(pRuterFunc func(pRouterInfo *mux.Router), pReadTimeout, pWriteTimeout, pIdleTimeout, pPortAdrs int) error {
	log.Info("SetServer (+)")

	if pReadTimeout == 0 {
		pReadTimeout = 30
	}
	if pWriteTimeout == 0 {
		pWriteTimeout = 30
	}
	if pIdleTimeout == 0 {
		pIdleTimeout = 120
	}

	lRouter := mux.NewRouter()
	pRuterFunc(lRouter)

	lRouter.MethodNotAllowedHandler = http.HandlerFunc(func(pResp http.ResponseWriter, pReq *http.Request) {
		if pReq.Method == http.MethodOptions {
			MsgSender(pResp, "Optional Call Success")
			return
		}
		pResp.Header().Set("Content-Type", "application/json")
		pResp.WriteHeader(http.StatusMethodNotAllowed)
		ErrorSender(pResp, "GMSS01", fmt.Errorf("Method %s not allowed on %s", pReq.Method, pReq.URL.Path))
	})

	// lHandler := logMiddleware(lRouter)
	lSrv := &http.Server{
		ReadTimeout:  time.Duration(pReadTimeout) * time.Second,
		WriteTimeout: time.Duration(pWriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(pIdleTimeout) * time.Second,
		Handler:      logMiddleware(lRouter),
		Addr:         fmt.Sprintf(":%d", pPortAdrs),
	}

	log.Info("server start on :%d ....", pPortAdrs)
	if len(osSignal) > 0 {
		go func() {
			if lErr := lSrv.ListenAndServe(); lErr != nil && lErr != http.ErrServerClosed {
				fmt.Println(time.Now(), lErr)
			}
		}()

		// Wait for SIGTERM / CTRL+C
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, osSignal...)
		<-sig
	} else {
		if lErr := lSrv.ListenAndServe(); lErr != nil {
			fmt.Println(time.Now(), lErr)
		}
	}

	log.Info("SetServer (-)")
	return nil
}
