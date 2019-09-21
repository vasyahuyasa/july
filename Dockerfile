FROM golang:1.13 as builder
WORKDIR /build
COPY . /build
RUN CGO_ENABLED=off go build ./cmd/july

FROM alpine:3.10.2
WORKDIR /app
COPY --from=builder /build/july /app/july
EXPOSE 80
CMD [ "/app/july" ]
