# Horus 𓂀

## Introduction

Horus is a request logger and viewer for [Go](https://golang.org). It allows developers log and view http requests made to their web application.

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

```go
...
if err = listener.Serve(":{preferred_port}", "{preferred_password}"); err != nil {
	log.Fatal(err)
}
...
```

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

You can explore the implementation in the [example folder](https://github.com/ichtrojan/horus/tree/main/example).

## Built by 

* Toni Akinmolayan - [twitter](https://twitter.com/toniastro_) [GitHub](https://github.com/toniastro)
* Michael Trojan Okoh - [twitter](https://twitter.com/ichtrojan) [GitHub](https://github.com/ichtrojan)



