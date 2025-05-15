FROM golang:1.21-alpine
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

# Build all binaries
RUN go build -o bin/rest ./cmd/rest
RUN go build -o bin/grpc ./cmd/grpc
RUN go build -o bin/graphql ./cmd/graphql

CMD ["./bin/rest"]