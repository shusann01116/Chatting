FROM golang:1.20.0 AS build

WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./...

FROM golang:1.20.0-alpine AS production

COPY --link --from=build /app .
CMD ["./main"]
