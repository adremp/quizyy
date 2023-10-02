FROM golang as build

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .
RUN go build -o main

EXPOSE 3000:3000

ENTRYPOINT ["/app/main"]