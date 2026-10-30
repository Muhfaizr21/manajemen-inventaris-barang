# syntax=docker/dockerfile:1

FROM golang:1.27-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/server .

FROM alpine:3.20
RUN adduser -D -u 10001 appuser
USER appuser
COPY --from=build /app/server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
