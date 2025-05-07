FROM golang:1.24-alpine

WORKDIR /app

RUN apk add --no-cache git sqlite

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server .

ENV SQLITE_FILENAME=books.db
ENV HOST=0.0.0.0
ENV PORT=3030

EXPOSE 3030

CMD ["./server"]