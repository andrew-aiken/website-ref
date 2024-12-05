FROM golang:1.23.3 AS build

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY src/ ./src

RUN CGO_ENABLED=0 go build -o /go/bin/crane-ecr-push ./src/main.go

# https://github.com/GoogleContainerTools/distroless
FROM gcr.io/distroless/static-debian12:nonroot-amd64

COPY --from=build /go/bin/crane-ecr-push /
ENTRYPOINT ["/crane-ecr-push"]
