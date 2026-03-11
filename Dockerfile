FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY prod/backend/go.mod prod/backend/go.sum ./
RUN go mod download

COPY prod/backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/server ./
COPY --from=builder /app/internal/model/words.go ./internal/model/
COPY --from=builder /app/induce_master.db ./

EXPOSE 8080

CMD ["./server"]
