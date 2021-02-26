package main

import (
	"fmt"
	"github.com/ichtrojan/horus"
	_ "github.com/ichtrojan/horus"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", horus.Watch(func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("https://google.com/")

		if err != nil {
			fmt.Print(err)
		}
		body , err := ioutil.ReadAll(resp.Body)

		defer resp.Body.Close()

		w.WriteHeader(resp.StatusCode)
		_, _ = fmt.Fprintf(w, string(body))
	}))

	if err := horus.Serve(":8016"); err != nil{
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":8888", nil); err != nil {
		fmt.Print(err)
	}
}
