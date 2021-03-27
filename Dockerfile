FROM golang:alpine AS build
WORKDIR /go/src/app
COPY . .
RUN go build -o /go/bin/app main.go

EXPOSE 8085

ENTRYPOINT ["/go/bin/app"]