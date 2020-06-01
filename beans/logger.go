package beans

import (
	"context"
	"log"
	"net/http"
)

// CtxKeyLogger 用来标识Context中存放的值。
const CtxKeyLogger key = "beans.Logger"

// Logger write to syslog
type Logger struct {
	*log.Logger
}

// NewLogger new Logger
func NewLogger(l *log.Logger) *Logger {
	l.Printf("logger init.")
	return &Logger{l}
}

// ServeHTTPimp implement pod.Handler
func (l *Logger) ServeHTTP(
	rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	ctx := context.WithValue(r.Context(), CtxKeyLogger, l.Logger)
	next(rw, r.WithContext(ctx))

}

// GetLogger return *log.Logger
func GetLogger(r *http.Request) (*log.Logger, bool) {
	l := r.Context().Value(CtxKeyLogger)

	if l == nil {
		return nil, false
	}

	return l.(*log.Logger), true
}
