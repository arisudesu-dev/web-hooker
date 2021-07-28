FROM golang:1.16 as builder
WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download
COPY . /app/
RUN CGO_ENABLED=0 go build -o web-hooker .

FROM alpine:3.14.0
WORKDIR /app
COPY --from=builder /app/web-hooker /app/
EXPOSE 8080
ENTRYPOINT ["/app/web-hooker"]
