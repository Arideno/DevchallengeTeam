FROM golang:1.15-alpine

WORKDIR /go/src/app
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .