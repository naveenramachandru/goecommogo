From golang:1.2.0

WORKDIR /app


COPY go.mod .
COPY main.go .

Run go build -o bin .

ENTRYPOINT ["/app/bin"]
