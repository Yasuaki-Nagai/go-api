FROM golang:1.20 AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o go-api .

FROM alpine:latest
COPY --from=build /app/go-api /go-api
CMD ["/go-api"]
