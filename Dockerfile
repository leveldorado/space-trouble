FROM golang:1.19-alpine AS builder

WORKDIR /space-trouble

COPY . .
RUN go mod tidy
# build app
RUN cd cmd/space-trouble && CGO_ENABLED=0 GOOS=linux go build

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /space-trouble/cmd/space-trouble/space-trouble .
CMD ["./space-trouble"]