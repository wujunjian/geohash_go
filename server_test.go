package main

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

//Testmain server
func TestMain(t *testing.T) {
	http.HandleFunc("/hello",
		func(w http.ResponseWriter, req *http.Request) {
			time.Sleep(1 * time.Second)
			fmt.Printf("receive msg form: %v", req.RemoteAddr)
			io.WriteString(w, "hello world!\n")

			time.Sleep(100 * time.Second)
			io.WriteString(w, "hello end!\n")

		})
	server := &http.Server{Addr: ":9990"}
	server.ListenAndServe()

}
