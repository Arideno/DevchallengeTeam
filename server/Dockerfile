FROM golang:1.15-alpine

RUN apk add build-base

WORKDIR /go/src/app
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .