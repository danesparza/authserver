package api

import (
	"fmt"
	"net/http"
)

// HelloWorld emits a hello world
func HelloWorld(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "Hello, world - service")
}
