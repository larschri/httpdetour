# HTTP Detour

Synchronous HTTP handling. Create `httpdetour.Chan` and inject it where `http.Handler` or `http.RoundTrip` is expected.
The `httpdetour.Chan` is a channel of request/response pairs that can be handled in a single goroutine. For example inside a test function.

## Install

    go get github.com/larschri/httpdetour

## Example

Create a `httpdetour.Chan` and inject it either as `http.Handler` or as `http.RoundTrip`.

```go
	handler := httpdetour.NewChan()
	server := httptest.NewServer(handler)

	r := <-handler
	defer r.Close()

	// Read r.Request
	// Write r.ResponseWriter
```
