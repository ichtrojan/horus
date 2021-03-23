package main

import (
	"encoding/json"
	"fmt"
	"github.com/ichtrojan/horus"
	"log"
	"net/http"
)

type response struct {
	Message string `json:"message"`
}

func main() {
	listener, err := horus.Init("mysql")

	if err != nil {
		log.Fatal(err)
	}
	if err = listener.Serve(":8081","12345"); err != nil{
		log.Fatal(err)
	}


	http.HandleFunc("/", listener.Watch(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		_ = json.NewEncoder(w).Encode(response{Message: "Horus is live 👁"})
	}))

	if err := http.ListenAndServe(":8888", nil); err != nil {
		fmt.Print(err)
	}
}
