# Build stage
FROM golang:1.24-alpine3.22 AS builer
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.22
WORKDIR /app
COPY --from=builer /app/main .

EXPOSE 1234
CMD [ "/app/main" ]