package horus

import (
	"fmt"
	"github.com/ichtrojan/horus/models"
	"github.com/ichtrojan/horus/storage"
	"net/http"
	"time"
)

func Watch(next func(http.ResponseWriter, *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := storage.Connect()

		if err != nil {
			_ = fmt.Errorf("%v", err)
		}
		fmt.Println("Logged req")
		req := models.Request{
			ResponseBody:  "",
			ResposeStatus: 200,
			RequestBody:   "",
			Path:          r.RequestURI,
			Headers:       "",
			Method:        r.Method,
			Host:          r.Host,
			Ipadress:      r.RemoteAddr,
			Time:          time.Now(),
		}

		write := request.Create(&req)

		fmt.Println(write)

		next(w, r)
	}
}

func Serve(port string) error {
	http.HandleFunc("/hor", func(w http.ResponseWriter, r *http.Request) {
		var req models.Request

		request, err := storage.Connect()

		if err != nil {
			_ = fmt.Errorf("%v", err)
		}

		request.First(&req)
	})

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		return err
	}

	return nil
}
