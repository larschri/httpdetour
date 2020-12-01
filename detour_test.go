package httpdetour

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func ExampleDetour() {
	detour := NewDetourChan()
	cli := http.Client{
		Transport: detour,
	}

	c := make(chan struct{})
	go func() {
		resp, _ := cli.Get("/ping")
		bs, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(bs))
		close(c)
	}()

	el := <-detour

	fmt.Println(el.Request.URL.Path)
	rec := httptest.NewRecorder()
	rec.Write([]byte("pong"))
	el.Respond(rec.Result(), nil)

	<-c
	// Output:
	// /ping
	// pong
}
