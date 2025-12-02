FROM golang:1.25.4 AS build
WORKDIR /go/src/app

COPY go.mod ./
RUN go mod download

COPY main.go ./main.go
RUN CGO_ENABLED=0 go build -o /go/bin/appName ./main.go

FROM alpine:latest

COPY --from=build /go/bin/appName /
ENTRYPOINT ["/appName"]
