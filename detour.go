package httpdetour

import (
	"net/http"
	"net/http/httptest"
)

// Exchange is the request/response pair that will be passed over a channel
type Exchange struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	done           chan<- bool
}

// Chan is a channel of request/response pairs. It can be injected where http.Handler or http.RoundTripper is expected.
type Chan chan *Exchange

// NewChan creates a new Chan
func NewChan() Chan {
	return make(chan *Exchange)
}

// ServeHTTP implements the http.Handler interface
func (r Chan) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	c := make(chan bool)
	elem := Exchange{req, res, c}

	select {
	case r <- &elem:
	case <-req.Context().Done():
		return
	}

	select {
	case <-c:
		return
	case <-req.Context().Done():
		return
	}
}

// RoundTrip implements the http.RoundTripper interface
func (r Chan) RoundTrip(req *http.Request) (*http.Response, error) {
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	return res.Result(), nil
}

// Close must be invoked to signal that the request has been handled
func (r *Exchange) Close() error {
	defer close(r.done)

	select {
	case r.done <- true:
		return nil
	case <-r.Request.Context().Done():
		return r.Request.Context().Err()
	}
}
