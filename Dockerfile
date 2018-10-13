FROM golang:1.11.1-alpine3.7

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
# RUN go install -v ./...

# CMD ["app"]
CMD ["go", "run", "main.go"]