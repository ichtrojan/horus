package horus

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/websocket"
	"github.com/ichtrojan/horus/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"html/template"
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
	key      string
}

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	filePeriod = 1 * time.Second
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	requestQueue = make(chan models.Request)
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

type Response struct {
	Message string
}

type Credentials struct {
	Key string
}

var tmpl = template.Must(template.ParseFiles("../views/auth.gohtml"))

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
		log.Fatal("HORUS_DB_PORT not set in .env")
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

		req := models.Request{
			ResponseBody:  string(minifyJson(recorder.Body.Bytes())),
			ResposeStatus: recorder.Code,
			RequestBody:   string(minifyJson(requestBody)),
			Path:          request.RequestURI,
			Headers:       string(minifyJson(headers)),
			Method:        request.Method,
			Host:          request.Host,
			Ipadress:      ipAddress,
			TimeSpent:     float64(time.Since(startTime)) / float64(time.Millisecond),
		}

		write := database.Create(&req)

		if write.RowsAffected != 1 {
			log.Fatal("unable to log request")
		}
		go func() {
			requestQueue <- req
		}()

		next(writer, request)
	}
}

func minifyJson(originalJson []byte) []byte {
	buffer := new(bytes.Buffer)

	if len(originalJson) == 0 {
		return []byte("[]")
	}

	if err := json.Compact(buffer, originalJson); err != nil {
		fmt.Println(err)
	}

	return []byte(buffer.String())
}

func (config Config) Serve(port string, key string) error {
	config.key = key

	horusServer := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("../views/public/"))

	horusServer.Handle("/public/", http.StripPrefix("/public", fileServer))

	horusServer.HandleFunc("/horus", renderView)

	horusServer.HandleFunc("/logs", config.showLogs)

	horusServer.HandleFunc("/ws", config.serveWs)

	horusServer.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			login(w, r)
		case "POST":
			config.postlogin(w, r)
		}
	})

	horusServer.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		cookie := &http.Cookie{
			Name:   "horus",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		}

		http.SetCookie(w, cookie)

		http.Redirect(w, r, "login", 302)
	})

	var err error

	go func() {
		err = http.ListenAndServe(port, horusServer)
	}()

	if err != nil {
		return err
	}

	return nil
}

func (config Config) postlogin(w http.ResponseWriter, r *http.Request) {

	creds := Credentials{
		Key: r.FormValue("key"),
	}

	if creds.Key == config.key {
		setSession(w, "god")

		response := map[string]bool{"status": true}

		_ = json.NewEncoder(w).Encode(response)

		return
	}

	response := map[string]bool{"status": false}

	_ = json.NewEncoder(w).Encode(response)

	return

}

func login(w http.ResponseWriter, r *http.Request) {
	_ = tmpl.Execute(w, nil)
}

func setSession(w http.ResponseWriter, who string) {
	value := map[string]string{
		"who": who,
	}

	if encoded, err := cookieHandler.Encode("horus", value); err == nil {
		cookie := &http.Cookie{
			Name:   "horus",
			Value:  encoded,
			Path:   "/",
			MaxAge: 3600,
		}

		http.SetCookie(w, cookie)
	}
}

func getSession(w http.ResponseWriter, r *http.Request) (who string) {
	if cookie, err := r.Cookie("horus"); err == nil {
		cookieValue := make(map[string]string)

		if err = cookieHandler.Decode("horus", cookie.Value, &cookieValue); err == nil {
			who = cookieValue["who"]
			return who
		}
	}
	return
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

func (config Config) showLogs(w http.ResponseWriter, r *http.Request) {
	lastID := r.URL.Query().Get("lastID")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var req []models.Request

	session := getSession(w, r)

	if session == ""{
		_ = json.NewEncoder(w).Encode(&req)
		return
	}

	request, err := connect(config)

	if err != nil {
		_ = fmt.Errorf("%v", err)
	}

	if lastID == "0" {
		request.Limit(20).Order("id desc").Find(&req)
	} else {
		request.Limit(20).Order("id desc").Where("id < ?", lastID).Find(&req)
	}

	_ = json.NewEncoder(w).Encode(&req)

	return
}

func renderView(w http.ResponseWriter, r *http.Request) {

	http.ServeFile(w, r, "../views/index.html")
}

func (config Config) serveWs(w http.ResponseWriter, r *http.Request) {

	session := getSession(w, r)

	if session == ""{

		response := map[string]string{"status": "Invalid session"}

		_ = json.NewEncoder(w).Encode(response)

		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}

	go writer(ws)

	reader(ws)
}

func reader(ws *websocket.Conn) {
	defer ws.Close()

	ws.SetReadLimit(512)

	_ = ws.SetReadDeadline(time.Now().Add(pongWait))

	ws.SetPongHandler(func(string) error {
		_ = ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := ws.ReadMessage()

		if err != nil {
			break
		}
	}
}

func writer(ws *websocket.Conn) {
	pingTicker := time.NewTicker(pingPeriod)

	defer func() {
		pingTicker.Stop()
		_ = ws.Close()
	}()

	for {
		select {
		case logs, ok := <-requestQueue:

			reqBodyBytes := new(bytes.Buffer)

			_ = json.NewEncoder(reqBodyBytes).Encode(logs)

			logsPush := reqBodyBytes.Bytes()

			_ = ws.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				_ = ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := ws.NextWriter(websocket.TextMessage)

			if err != nil {
				return
			}

			_, _ = w.Write(logsPush)

			n := len(requestQueue)

			for i := 0; i < n; i++ {
				_ = json.NewEncoder(reqBodyBytes).Encode(requestQueue)
				_, _ = w.Write(reqBodyBytes.Bytes())
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-pingTicker.C:
			_ = ws.SetWriteDeadline(time.Now().Add(writeWait))

			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
