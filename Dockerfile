FROM golang:1.17-alpine3.14

WORKDIR /app
ADD *.go templates ./
ADD go.mod go.sum ./
RUN go get -v
RUN CGO_ENABLED=0 GOOS=linux go test -v ./...
RUN go build -o /app/bin/zcdash

CMD ["/app/bin/zcdash"]
