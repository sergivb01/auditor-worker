# build stage
FROM golang:1.14.0-alpine as builder
ENV GO111MODULE=on
WORKDIR /app
# copy go.mod and sum first for caching
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o application .


# final stage
FROM scratch
COPY --from=builder /app/application /app/

# port configuration
EXPOSE 8080
ENV LISTEN_ADDR ":8080"

USER default

ENTRYPOINT ["/app/application"]