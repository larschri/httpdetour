package httpdetour

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
)

func ExampleDetour() {
	detour := NewChan()
	server := httptest.NewServer(detour)

	go func() {
		d := <-detour
		defer d.Close()

		d.ResponseWriter.Write([]byte("hello world"))
	}()

	resp, _ := http.Get(server.URL)
	io.Copy(os.Stdout, resp.Body)

	// Output: hello world
}

func ExampleDetourRoundTrip() {
	detour := NewChan()
	cli := http.Client{
		Transport: detour,
	}

	go cli.Get("/hello-world")

	el := <-detour

	fmt.Println(el.Request.URL.String())
	el.Close()

	// Output: /hello-world
}
