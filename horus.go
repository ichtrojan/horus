package horus

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ichtrojan/horus/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"time"
)

type Config struct {
	Database string
	Dsn      string
}

func Init(database string) (Config, error) {
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
		log.Fatal("HORUS_DB_NAME not set in .env")
	}

	switch database {
	case "mysql":
		return Config{
			Database: "mysql",
			Dsn:      fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, name),
		}, nil
	case "postgres":
		return Config{
			Database: "mysql",
			Dsn:      fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", host, user, pass, name, port),
		}, nil
	default:
		return Config{}, errors.New("database not defined")
	}
}

func (config Config) Watch(next func(http.ResponseWriter, *http.Request)) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		database, err := connect(config)

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

		fmt.Println(recorder.Body.String())

		req := models.Request{
			ResponseBody:  recorder.Body.String(),
			ResposeStatus: recorder.Code,
			RequestBody:   string(requestBody),
			Path:          request.RequestURI,
			Headers:       string(headers),
			Method:        request.Method,
			Host:          request.Host,
			Ipadress:      ipAddress,
			TimeSpent:     float64(time.Since(startTime)) / float64(time.Millisecond),
		}

		write := database.Create(&req)

		if write.RowsAffected != 1 {
			log.Fatal("unable to log request")
		}

		next(writer, request)
	}
}

func (config Config) Serve(port string) error {
	http.HandleFunc("/horus", func(w http.ResponseWriter, r *http.Request) {
		var req models.Request

		request, err := connect(config)

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

func connect(config Config) (*gorm.DB, error) {
	db, err := gorm.Open(config.Database, config.Dsn)

	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.Request{}).Error; err != nil {
		return nil, err
	}

	return db, nil
}
