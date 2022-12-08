FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/hermes ./main.go

CMD ["hermes"]