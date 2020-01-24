FROM golang:1.13 as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN make



FROM debian:stable

WORKDIR /app

COPY --from=builder /app/cmd/build/* ./

EXPOSE 80

CMD ["./calendar_sender"]
