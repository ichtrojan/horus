package horus

import (
	"encoding/json"
	"fmt"
	"github.com/ichtrojan/horus/models"
	"github.com/ichtrojan/horus/storage"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"time"
)

func Watch(next func(http.ResponseWriter, *http.Request)) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		database, err := storage.Connect()

		if err != nil {
			log.Fatal(err)
		}

		ipAddress, _, _ := net.SplitHostPort(request.RemoteAddr)

		headers, err := json.Marshal(request.Header)

		if err != nil {
			log.Fatal(err)
		}

		requestBody, err := ioutil.ReadAll(request.Body)

		if err != nil {
			log.Fatal(err)
		}

		recorder := httptest.NewRecorder()

		startTime := time.Now()

		next(recorder, request)

		req := models.Request{
			ResponseBody:  recorder.Body.String(),
			ResposeStatus: recorder.Code,
			RequestBody:   requestBody,
			Path:          request.RequestURI,
			Headers:       headers,
			Method:        request.Method,
			Host:          request.Host,
			Ipadress:      ipAddress,
			TimeSpent:     float64(time.Since(startTime)) / float64(time.Millisecond),
		}

		write := database.Create(&req)

		if write.RowsAffected != 1 {
			log.Fatal("unable to log request")
		}
	}
}

func Serve(port string) error {
	http.HandleFunc("/horus", func(w http.ResponseWriter, r *http.Request) {
		var req models.Request

		request, err := storage.Connect()

		if err != nil {
			_ = fmt.Errorf("%v", err)
		}

		request.First(&req)

	})

	go func() error {
		if err := http.ListenAndServe(port, nil); err != nil {
			return err
		}
		return nil
	}()

	fmt.Println("Started horus:views server on port" + port)

	return nil
}
