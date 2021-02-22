package horus

import (
	"fmt"
	"github.com/ichtrojan/horus/models"
	"github.com/ichtrojan/horus/storage"
	"net/http"
)

func Watch(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request, err := storage.Connect()

		if err != nil {
			_ = fmt.Errorf("%v", err)
		}

		request.Create(models.Request{
			ResponseBody:  "",
			ResposeStatus: r.Response.StatusCode,
			RequestBody:   "",
			Path:          r.RequestURI,
			Headers:       "",
			Method:        r.Method,
			Host:          r.Host,
		})

		next.ServeHTTP(w, r)
	})
}

func Serve(port string) error {
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		return err
	}

	return nil
}
