FROM golang:1.20 AS builder

WORKDIR /build

RUN go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY Makefile go.mod go.sum ./
RUN make download

COPY . .
RUN make keeper

FROM debian:bookworm-slim

WORKDIR /app

RUN mkdir ssl

COPY --from=builder /go/bin/migrate .
COPY --from=builder /build/ssl/ca/keeper.crt ./ssl/
COPY --from=builder /build/ssl/ca/keeper.key ./ssl/
COPY --from=builder /build/scripts/init-keeper .
COPY --from=builder /build/migrations ./migrations
COPY --from=builder /build/dist/keeper .

RUN chown -R nobody:nogroup /app \
    && chmod -R 777 /app

USER nobody

ENTRYPOINT ["./keeper"]
