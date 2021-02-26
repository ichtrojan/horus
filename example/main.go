package main

import (
	"encoding/json"
	"fmt"
	"github.com/ichtrojan/horus"
	_ "github.com/ichtrojan/horus"
	"net/http"
)

type response struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/", horus.Watch(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		_ = json.NewEncoder(w).Encode(response{Message: "Horus is live 👁"})
	}))

	if err := http.ListenAndServe(":8888", nil); err != nil {
		fmt.Print(err)
	}
}
