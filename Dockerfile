#build stage
FROM golang:alpine AS builder
WORKDIR /go/src/app
COPY . .
RUN apk add --no-cache git
RUN go get -d -v ./...
RUN go install -v ./...
RUN ls /go/bin/

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/blog-api /app
ENTRYPOINT ./app
LABEL Name=go-api Version=0.0.1
EXPOSE 8080
