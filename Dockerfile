FROM golang:1.21-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o ./go-clean-architecture ./cmd

FROM alpine
COPY --from=build /src/go-clean-architecture /usr/local/bin/app
COPY --from=build /src/config /config

ENV GO_CLEAN_ARCHITECTURE_API_HTTP_HOST=0.0.0.0

ENTRYPOINT ["app", "-c", "/config/config.yaml"]

EXPOSE 8080