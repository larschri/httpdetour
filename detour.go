package httpdetour

import (
	"net/http"
)

type DetourResponse struct {
	Response *http.Response
	Error    error
}

type Detour struct {
	Request      *http.Request
	ResponseChan chan<- DetourResponse
}

func (r *Detour) Respond(resp *http.Response, err error) error {
	defer close(r.ResponseChan)

	select {
	case r.ResponseChan <- DetourResponse{resp, err}:
		return nil
	case <-r.Request.Context().Done():
		return r.Request.Context().Err()
	}
}

type DetourChan chan *Detour

func NewDetourChan() DetourChan {
	return make(chan *Detour)
}

func (r DetourChan) RoundTrip(req *http.Request) (*http.Response, error) {
	c := make(chan DetourResponse)
	elem := Detour{req, c}

	select {
	case r <- &elem:
	case <-req.Context().Done():
		return nil, req.Context().Err()
	}

	select {
	case resp := <-c:
		return resp.Response, resp.Error
	case <-req.Context().Done():
		return nil, req.Context().Err()
	}
}
