package main

import (
	"fmt"
	"github.com/ichtrojan/horus"
	_ "github.com/ichtrojan/horus"
	"net/http"
)

func main() {
	http.HandleFunc("/", horus.Watch(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "Hello, Testing from Thoth")
	}))

	//_ = horus.Serve("8000")

	if err := http.ListenAndServe(":8888", nil); err != nil {
		fmt.Print(err)
	}
}
