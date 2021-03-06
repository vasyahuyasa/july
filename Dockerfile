FROM golang:1.14.2-alpine3.11 as builder
WORKDIR /build
COPY . /build
RUN CGO_ENABLED=off go build ./cmd/july

FROM alpine:3.11
WORKDIR /app
COPY --from=builder /build/july /app/july
EXPOSE 80
ENV STORAGE_DRIVER=gdrive
ENV CATALOG_ROOT=root
ENV GOOGLE_CREDENTIALS_PATH=credentials.json
env GOOGLE_TOKEN_PATH=token.json
CMD /app/july \
    -s $STORAGE_DRIVER \
    -d $CATALOG_ROOT \
    -googlecred $GOOGLE_CREDENTIALS_PATH \
    -googletoken $GOOGLE_TOKEN_PATH
