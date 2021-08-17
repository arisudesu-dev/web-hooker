FROM golang:1.17.0 as builder
WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download
COPY . /app/
RUN CGO_ENABLED=0 go build -o web-hooker .

FROM alpine:3.14.1
WORKDIR /app/scripts
COPY --from=builder /app/web-hooker /app/
EXPOSE 8080
ENTRYPOINT ["/app/web-hooker"]
