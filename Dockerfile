FROM golang:1.24-alpine AS builder
WORKDIR /app


COPY ./ ./
RUN go mod download

RUN go build -o /lighthouse github.com/go-oidfed/lighthouse
RUN go build -o /lhcli github.com/go-oidfed/lighthouse

FROM debian:stable
RUN apt-get update && apt-get install -y ca-certificates && apt-get autoremove -y && apt-get clean -y && rm -rf /var/lib/apt/lists/*

COPY --from=builder /lighthouse .
COPY --from=builder /lhcli .

CMD bash -c "update-ca-certificates && /lighthouse /data/config.yaml"
