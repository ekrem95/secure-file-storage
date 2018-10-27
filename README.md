## secure-file-storage

The stored image file is secured, as the file is being encrypted not by just using one but two encryption algorithm which are AES and DES. Data is kept secured on server which avoids unauthorized access.

#### Prerequisites

```
go1.11.1 or higher
A working golang environment
Docker/Docker Compose
```

#### Installing

```
go get github.com/ekrem95/secure-file-storage
```

Go to project directory and run `docker-compose.yml` file to start databases

```
cd $GOPATH/src/github.com/ekrem95/secure-file-storage && docker-compose up -d
```

Run the tests

```
go test -v ./...
```

Run the app with `go run main.go` or `secure-file-storage`
