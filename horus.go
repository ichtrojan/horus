package horus

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/websocket"
	"github.com/ichtrojan/horus/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var METHODS = map[string]string{
	"get":    "GET",
	"post":   "POST",
	"delete": "DELETE",
	"option": "OPTION",
	"put":    "PUT",
}

type InternalConfig struct {
	Database string
	Dsn      string
	key      string
	db       *gorm.DB
}

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
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

type Config struct {
	DbUser    string
	DbPssword string
	DbHost    string
	DbName    string
	DbPort    string
}

func Init(database string, config Config) (*InternalConfig, error) {
	var dsn string

	user := config.DbUser

	pass := config.DbPssword

	host := config.DbHost

	name := config.DbName

	port := config.DbPort

	switch database {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, name)
	case "postgres":
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", host, user, pass, name, port)
	default:
		msg := fmt.Sprintf("Database %s is not supported by horus", database)
		return nil, errors.New(msg)
	}

	db, err := connect(database, dsn)
	if err != nil {
		return nil, err
	}

	return &InternalConfig{
		Database: database,
		Dsn:      dsn,
		db:       db,
	}, nil
}

func (config *InternalConfig) Watch(next func(http.ResponseWriter, *http.Request)) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		headers, err := json.Marshal(request.Header)

		if err != nil {
			fmt.Println(err)
		}

		requestBody, err := ioutil.ReadAll(request.Body)

		if err != nil {
			fmt.Println(err)
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
			Ipadress:      getIp(request),
			TimeSpent:     float64(time.Since(startTime)) / float64(time.Millisecond),
		}

		write := config.db.Create(&req)

		if write.RowsAffected != 1 {
			fmt.Println("unable to log request")
		}

		request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))

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

func (config *InternalConfig) Serve(port string, key string) error {
	config.key = key

	_, filename, _, _ := runtime.Caller(0)

	horusServer := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(path.Dir(filename) + "/views/public/"))

	horusServer.Handle("/public/", http.StripPrefix("/public", fileServer))

	horusServer.HandleFunc("/horus", renderView)

	horusServer.HandleFunc("/logs", config.showLogs)

	horusServer.HandleFunc("/ws", config.serveWs)

	horusServer.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
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

		http.Redirect(w, r, "horus", 302)
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

func (config *InternalConfig) postlogin(w http.ResponseWriter, r *http.Request) {
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

func connect(database string, dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(database, dsn)

	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.Request{}).Error; err != nil {
		return nil, err
	}

	return db, nil
}

func (config *InternalConfig) showLogs(w http.ResponseWriter, r *http.Request) {
	lastID := r.URL.Query().Get("lastID")
	method := r.URL.Query().Get("method")
	method = METHODS[method]

	var req []models.Request

	session := getSession(w, r)

	if session == "" {
		_ = json.NewEncoder(w).Encode(&req)
		return
	}

	if method == "" {
		method = "%"
	}

	if lastID == "0" {
		config.db.Limit(20).Order("id desc").Where("method LIKE ?", method).Find(&req)
	} else {
		config.db.Limit(20).Order("id desc").Where("id < ? AND method LIKE ?", lastID, method).Find(&req)
	}

	_ = json.NewEncoder(w).Encode(&req)

	return
}

func renderView(w http.ResponseWriter, r *http.Request) {
	_, filename, _, _ := runtime.Caller(0)

	http.ServeFile(w, r, path.Dir(filename)+"/views/index.html")
}

func (config *InternalConfig) serveWs(w http.ResponseWriter, r *http.Request) {
	session := getSession(w, r)

	if session == "" {
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

func (config *InternalConfig) Close() error {
	return config.db.DB().Close()
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

func getIp(request *http.Request) string {
	var xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
	var xRealIP = http.CanonicalHeaderKey("X-Real-IP")

	var ip string

	if xff := request.Header.Get(xForwardedFor); xff != "" {
		i := strings.Index(xff, ", ")

		if i == -1 {
			i = len(xff)
		}

		ip = xff[:i]

	} else if xrip := request.Header.Get(xRealIP); xrip != "" {
		ip = xrip
	}

	if ip == "" {
		ipAddress, _, err := net.SplitHostPort(request.RemoteAddr)

		if err != nil {
			return ""
		}

		return strings.TrimSpace(ipAddress)
	}

	return ip
}
