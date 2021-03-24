package main

import (
	"encoding/json"
	"fmt"
	"github.com/ichtrojan/horus"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	user, exist := os.LookupEnv("HORUS_DB_USER")

	if !exist {
		log.Fatal("HORUS_DB_USER not set in .env")
	}

	pass, exist := os.LookupEnv("HORUS_DB_PASS")

	if !exist {
		log.Fatal("HORUS_DB_PASS not set in .env")
	}

	host, exist := os.LookupEnv("HORUS_DB_HOST")

	if !exist {
		log.Fatal("HORUS_DB_HOST not set in .env")
	}

	name, exist := os.LookupEnv("HORUS_DB_NAME")

	if !exist {
		log.Fatal("HORUS_DB_NAME not set in .env")
	}

	port, exist := os.LookupEnv("HORUS_DB_PORT")

	if !exist {
		log.Fatal("HORUS_DB_PORT not set in .env")
	}

	listener, err := horus.Init("mysql", horus.Config{
		DbName:    name,
		DbHost:    host,
		DbPssword: pass,
		DbPort:    port,
		DbUser:    user,
	})

	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":8081", listener.Serve("12345")); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", listener.Watch(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)

			response := map[string]string{"message": "endpont not found"}

			_ = json.NewEncoder(w).Encode(response)

			return
		}

		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)

			response := map[string]string{"message": "method not allowed"}

			_ = json.NewEncoder(w).Encode(response)

			return
		}

		response := map[string]string{"message": "Horus is live 👁"}

		_ = json.NewEncoder(w).Encode(response)
	}))

	http.HandleFunc("/message", listener.Watch(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Path != "/message" {
			w.WriteHeader(http.StatusNotFound)

			response := map[string]string{"message": "endpont not found"}

			_ = json.NewEncoder(w).Encode(response)

			return
		}

		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)

			response := map[string]string{"message": "method not allowed"}

			_ = json.NewEncoder(w).Encode(response)

			return
		}

		response := map[string]string{"message": "message received"}

		_ = json.NewEncoder(w).Encode(response)
	}))

	if err := http.ListenAndServe(":8888", nil); err != nil {
		fmt.Print(err)
	}
}
