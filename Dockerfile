FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod ./
# Dependencies will be downloaded when go mod download runs; go.sum is generated on first tidy.
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 1323
CMD ["./server"]
