FROM golang:1.9 as builder

WORKDIR /go/src/todo-backend
COPY . .

RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -installsuffix netgo -ldflags '-w' -o main ./cmd/todo

FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
MAINTAINER Jose Ordiales <jlordiales@gmail.com>
EXPOSE 8080

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/todo-backend/main /main

ENTRYPOINT ["/main"]
