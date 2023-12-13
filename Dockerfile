FROM golang:1.21.1-alpine AS builder

RUN wget -O /usr/local/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.5/dumb-init_1.2.5_x86_64
RUN chmod +x /usr/local/bin/dumb-init

WORKDIR /app
COPY . /app/

RUN go mod tidy
RUN go build -o tg-tefcon-2023 main.go

FROM gcr.io/distroless/base-debian11
COPY --from=builder /app/tg-tefcon-2023 /bin/
COPY --from=builder /usr/local/bin/dumb-init /bin/
USER nonroot
ENTRYPOINT ["/bin/dumb-init", "--", "/bin/tg-tefcon-2023"]
