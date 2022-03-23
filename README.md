# Horus 𓂀

<img width="1461" alt="horus-hero" src="https://user-images.githubusercontent.com/5338836/112303700-e2a23480-8c9c-11eb-9e47-9b0634e5b1e5.png">

## Introduction

Horus is a request logger and viewer for [Go](https://golang.org). It allows developers log and view http requests made to their web application.

<img width="1461" alt="horus-big-brother" src="https://user-images.githubusercontent.com/5338836/112304004-44629e80-8c9d-11eb-930b-bdf32448673c.png">

## Installation 

Run the following command to install Horus on your project:

```bash
go get github.com/ichtrojan/horus
```

### Initiate horus

```go
package main

import github.com/ichtrojan/horus

func main() {
    listener, err := horus.Init("mysql", horus.Config{
		DbName:    "{preferred_database_name}",
		DbHost:    "{preferred_database_host}",
		DbPssword: "{preferred_database_password}",
		DbPort:    "{preferred_database_port}",
		DbUser:    "{preferred_database_user}",
	})
}
```

>**NOTE**<br/>
> Supported database adapters include **mysql** and **postgres**

### Serve dashboard (optional)

<img width="1461" alt="horus-message-received" src="https://user-images.githubusercontent.com/5338836/112304136-69571180-8c9d-11eb-80db-167c6e8c4d3c.png">

```go
...
if err := listener.Serve(":8081", "{preferred_password}")); err != nil {
	log.Fatal(err)
}
...
```

>**NOTE**<br/>
> Visit `/horus` on the port configured to view the dashboard

### Usage
To enable horus to listen for requests, use the `Watch` middleware provided by horus on the endpoints you will like monitor.

```go
...
http.HandleFunc("/", listener.Watch(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    response := map[string]string{"message": "Horus is live 👁"}

    _ = json.NewEncoder(w).Encode(response)
}))
...
```

Remember to either `defer` or manually close the horos listener when you get a cancel signal on your server
```go
defer listener.Close()
```

or 

```go
signal.Notify(sigChan, os.Kill)
signal.Notify(sigChan, os.Interrupt)

signal := <-sigChan

log.Printf("Received %s signal, gracefully shutting down", signal)
if err := listener.Close(); err != nil {
	log.Fatalf("FATAL: Error while shutting down horus: %s", err)
}
```

You can explore the implementation in the [example folder](https://github.com/ichtrojan/horus/tree/main/example).

## Built by 

* Toni Akinmolayan - [twitter](https://twitter.com/toniastro_) [GitHub](https://github.com/toniastro)
* Michael Trojan Okoh - [twitter](https://twitter.com/ichtrojan) [GitHub](https://github.com/ichtrojan)
