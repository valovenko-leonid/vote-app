# stage 1 ─ build
FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

# stage 2 ─ small runtime
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/server .
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["./server"]