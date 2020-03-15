# build stage
FROM golang:1.14.0-alpine as builder
ENV GO111MODULE=on
# copy go.mod and sum first for caching
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
# GOOS=linux GOARCH=amd64
RUN CGO_ENABLED=0 go build -o /app/application .


# final stage
FROM scratch
COPY --from=builder /app/application /app/

# port configuration
EXPOSE 8080
ENV LISTEN_ADDR ":8080"

USER default

ENTRYPOINT ["/app/application"]