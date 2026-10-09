FROM golang:1.26.2-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/cmd/booking/exe /app/cmd/booking


FROM alpine:3.23
WORKDIR /app
COPY --from=builder /app/cmd/booking/exe /app
ENV PROJECT_ROOT=/app
CMD ["/app/exe"]