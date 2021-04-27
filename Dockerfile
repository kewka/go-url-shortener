FROM golang:1.16.4-alpine3.12 as builder
ENV CGO_ENABLED 0
RUN apk add make curl
RUN curl -L https://github.com/ufoscout/docker-compose-wait/releases/download/2.9.0/wait -o /wait && chmod +x /wait
WORKDIR /app
COPY . .
RUN make build

FROM busybox:1.33.1-musl
COPY --from=builder /app/bin /bin
COPY --from=builder /wait /wait
CMD /wait && /bin/go-url-shortener