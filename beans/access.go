package beans

import (
	"log"
	"net/http"
	"time"
)

// responseWriter warp http.ResponseWriter.
// There is no way to get http status code and written size by
// using default http.ResponseWriter.
type responseWriter struct {
	http.ResponseWriter
	wroteHeader bool
	status      int
	size        int
}

// WriteHeader warp http.ResponseWriter.WriteHeader
func (rw *responseWriter) WriteHeader(s int) {
	if !rw.wroteHeader {
		rw.wroteHeader = true
		rw.status = s
	}

	rw.ResponseWriter.WriteHeader(s)
}

// Write warp http.ResponseWriter.Write
func (rw *responseWriter) Write(b []byte) (int, error) {

	// look at http.ResponseWriter.Write() implementation
	// 虽然 ResponseWriter.Write() 会保底设置 http.StatusOK，
	// 但是只能调用内部的 ResponseWriter.WriteHeader？导致这里无法
	// 拿到保底状态码 http.StatusOK，日志里会出现HTTP返回码为0的情况。
	// 所以主动调用一下重写的 WriteHeader，修正之。
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}

	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

// Size of http server written bytes
func (rw *responseWriter) Size() int {
	return rw.size
}

// Status return http server status code
func (rw *responseWriter) Status() int {
	return rw.status
}

// Access write access log with log.Logger
type Access struct {
	*log.Logger
}

// NewAccess new Access
func NewAccess(l *log.Logger) *Access {
	return &Access{l}
}

// ServeHTTP implement pod.Handler
func (a *Access) ServeHTTP(
	rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	timeStart := time.Now()

	newrw := &responseWriter{rw, false, 0, 0}
	httpMethod := r.Method
	urlPath := r.URL.String()

	next(newrw, r)

	timeEnd := time.Now()
	du := timeEnd.Sub(timeStart)

	// TODO: 增加自定义日志输出格式？
	a.Printf(
		"%s %s %s %s %d %d",
		r.RemoteAddr, httpMethod, urlPath, du, newrw.Size(), newrw.Status())
}
