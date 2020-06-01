package pod

import (
	"container/list"
	"log"
	"net/http"
)

// Handler interface
// http.HandleFunc is the entry of next Hander
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request, http.HandlerFunc)
}

// Pod is tiny web framwork.
type Pod struct {
	handlers *list.List
	handle   http.HandlerFunc
}

// New Pod
func New() *Pod {
	p := &Pod{list.New(), nil}
	p.rebuild()
	return p
}

// rebuild handlers stack when any pushed in
func (p *Pod) rebuild() {
	p.handle = chainHandler(p.handlers.Front())
}

// ServeHTTP implement http.Handler
func (p *Pod) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	p.handle(rw, r)
}

// Push one Handler to call stack
func (p *Pod) Push(h Handler) {
	p.handlers.PushBack(h)
	p.rebuild()
}

// Run bind address and serve in http
func (p *Pod) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, p))
}

// RunTLS bind address and serve in https
func (p *Pod) RunTLS(addr, certFile, keyFile string) {
	log.Fatal(http.ListenAndServeTLS(addr, certFile, keyFile, p))
}

// chainHandler push handler to stack
func chainHandler(el *list.Element) http.HandlerFunc {
	// last element of the handlers chain is nil.
	if el == nil {
		return func(rw http.ResponseWriter, r *http.Request) {}
	}
	return func(rw http.ResponseWriter, r *http.Request) {
		el.Value.(Handler).ServeHTTP(rw, r, chainHandler(el.Next()))
	}
}
