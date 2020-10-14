FROM golang:alpine as builder
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ARG TARGETOS
ARG TARGETARCH
WORKDIR /app
COPY . .
RUN go build

FROM alpine:edge
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/garage-go /app/garage-go
EXPOSE 8080
ENTRYPOINT [ "sh", "-c", "/app/garage-go" ]

