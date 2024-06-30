FROM golang:1.20.5-alpine as intialize
RUN mkdir /app


WORKDIR /app
ADD . /app

RUN mkdir /app/bin

FROM intialize as builder
RUN go mod download
RUN go build -o ./bin  ./...


FROM gcr.io/distroless/base
COPY --from=builder /app/bin /app/bin

WORKDIR /app

ENV PORT 3000
EXPOSE 3000

CMD ["/app/bin/zephyrian", "service", "web"]
