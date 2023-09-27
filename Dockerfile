FROM golang

WORKDIR /app

COPY . .

RUN go build -o ./main

EXPOSE 3000:3000

ENTRYPOINT ["./main"]