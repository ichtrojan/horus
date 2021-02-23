package horus

import (
	"fmt"
	"github.com/ichtrojan/horus/models"
	"github.com/ichtrojan/horus/storage"
	"io"
	"net/http"
	"github.com/gorilla/handlers"
	"time"
)

func Watch(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := storage.Connect()

		if err != nil {
			_ = fmt.Errorf("%v", err)
		}
		fmt.Println("Logged")
		storage.GetDB().Create(&models.Request{
			ResponseBody:  "",
			ResposeStatus: 100,
			RequestBody:   "",
			Path:          r.RequestURI,
			Headers:       "",
			Method:        r.Method,
			Host:          r.Host,
			Ipadress:      r.RemoteAddr,
			Time: 		   time.Now(),
		})

		next.ServeHTTP(w, r)
	})
}

func newLoggingHandler(dst io.Writer) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(dst, h)
	}
}

func Serve(port string) error {
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		return err
	}

	return nil
}
