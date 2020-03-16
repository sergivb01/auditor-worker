# build stage
FROM golang:1.14.0-buster as builder

# if dependencies are updated, rebuild everything
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go mod verify

COPY . .
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o application github.com/sergivb01/auditor-worker/cmd

#second stage
FROM scratch
WORKDIR /application/

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/application .

CMD ["./application"]

#export DOCKER_BUILDKIT=1