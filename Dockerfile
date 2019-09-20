FROM golang:1.13 as builder
WORKDIR /build
COPY . /build
RUN go build ./cmd/july

FROM alpine
WORKDIR /app
COPY --from=builder /build/july /app/july
EXPOSE 80
CMD [ "/app/july" ]
