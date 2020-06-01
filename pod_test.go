package pod

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type testHandler struct {
	count int
}

func (th *testHandler) ServeHTTP(
	rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	th.count++
	next(rw, r)
}

// test pod
func Test_New_pod(t *testing.T) {
	p := New()
	if reflect.TypeOf(p.handle).String() != "http.HandlerFunc" {
		t.Errorf("pod.handle shoult be type of http.HandlerFunc.")
	}
	if p.handlers.Len() != 0 {
		t.Errorf("pod.handlers after init is not empty.")
	}
	t.Log("PASSED")
}

func Test_pod_ServeHTTP(t *testing.T) {
	rw := httptest.NewRecorder()
	r, err := http.NewRequest(
		"GET", "http://localhost/soLongDefiniteNotExist", nil)
	if err != nil {
		t.Error(err)
	}

	p := New()
	p.ServeHTTP(rw, r)
	t.Log("PASSED")
}

func Test_pod_Use(t *testing.T) {
	rw := httptest.NewRecorder()
	r, err := http.NewRequest(
		"GET", "http://localhost/soLongDefiniteNotExist", nil)
	if err != nil {
		t.Error(err)
	}

	th := &testHandler{count: 0}
	p := New()
	p.Push(th)

	loops := 100
	for i := 0; i < loops; i++ {
		p.ServeHTTP(rw, r)
	}

	if th.count != 100 {
		t.Errorf("handler executed more than expected. %d != %d",
			loops, th.count)
	}
	t.Log("PASSED")
}

func Test_pod_Run(t *testing.T) {
	t.Log("PASSED")
}
