package main

import (
	"fmt"
	"log"
	"log/syslog"
	"net/http"

	"github.com/shuaiming/pod"
	"github.com/shuaiming/pod/beans"
	"github.com/shuaiming/pod/beans/sessions"
)

func main() {

	l, _ := syslog.NewLogger(syslog.LOG_INFO|syslog.LOG_LOCAL0, log.Lmsgprefix)

	access := beans.NewAccess(l)
	static := beans.NewStatic("/public/", http.Dir("/tmp"), true)
	store := sessions.NewFilesystemStore(3600, "/tmp/.sessions")
	_sessions := sessions.New(store, 3600, 30, "id")

	// _sessions := sessions.New(sessions.NewMemoryStore(3600), 3600, 30, "id")
	_openid := beans.NewOpenID(
		"/openid",
		"http://you.domain.com:8000",
		"https://login.provider.com/openid/")

	logger := beans.NewLogger(l)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {

		logger, _ := beans.GetLogger(r)
		logger.Printf("%s %s", r.Method, r.URL.Path)

		if r.URL.Path != "/" {
			http.NotFound(rw, r)
			return
		}

		session, ok := sessions.GetSession(r)
		if !ok {
			http.Error(rw, "500 InternalServerError", http.StatusInternalServerError)
			return
		}

		// json.Unmarshal default use int64
		count, ok := session.Load("count")
		if !ok {
			count = 1
		} else {
			count = count.(int) + 1
		}

		fmt.Fprintf(rw, "Wherecome, the count is %d", count)
		session.Store("count", count)
		fmt.Fprintf(rw, "\n")

		if user, ok := beans.GetOpenIDUser(session); ok {
			for key, value := range user {
				fmt.Fprintf(rw, "%s -> %s\n", key, value)
			}
		}
	})

	app := pod.New()
	app.Push(access)
	app.Push(logger)
	app.Push(static)
	app.Push(_sessions)
	app.Push(_openid)
	app.Push(beans.NewHandler(mux))

	app.Run(":8000")
}
