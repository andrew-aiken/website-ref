FROM public.ecr.aws/docker/library/golang:1.24.3-alpine3.21 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ /app/cmd
COPY pkg/ /app/pkg
RUN CGO_ENABLED=0 go build -o /cloudfront-frontend /app/cmd/main.go


FROM public.ecr.aws/docker/library/node:24.1.0-alpine3.21 AS production
USER node

COPY --chown=root:node --chmod=010 --from=build /cloudfront-frontend /cloudfront-frontend

ENTRYPOINT [ "/cloudfront-frontend" ]
