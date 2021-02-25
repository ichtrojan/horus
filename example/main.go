package main

import (
	"fmt"
	"github.com/ichtrojan/horus"
	_ "github.com/ichtrojan/horus"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", horus.Watch(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		_, _ = fmt.Fprintf(w, "Hello, Testing from Thoth")
	}))

	if err := horus.Serve(":8016"); err != nil{
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":8888", nil); err != nil {
		fmt.Print(err)
	}
}
