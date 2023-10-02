FROM golang:alpine as build

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .
RUN go build -o main

FROM golang:alpine
COPY --from=build /app/main .
COPY --from=build /app/.env .
COPY --from=build /app/templates ./templates
COPY --from=build /app/static ./static

EXPOSE 3000

ENTRYPOINT [ "./main" ]