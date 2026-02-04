FROM golang:1.22-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/dashboard ./

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /out/dashboard ./dashboard
COPY --from=builder /src/assets ./assets
ENV PORT=9002 COUNTING_SERVICE_URL=http://counting:9003
EXPOSE 9002
ENTRYPOINT ["./dashboard"]