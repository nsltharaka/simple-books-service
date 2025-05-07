FROM golang:1.24-alpine

# Enable CGO and install C toolchain
ENV CGO_ENABLED=1
ENV GOOS=linux

RUN apk add --no-cache git sqlite build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server .

ENV SQLITE_FILENAME=books.db
ENV HOST=0.0.0.0
ENV PORT=3030

EXPOSE 3030

CMD ["./server"]