FROM golang:1.21.1-alpine3.18 as builder

WORKDIR /

COPY go.* ./
RUN go mod download 

COPY . . 
COPY dev.env .
RUN go build -o /myapp .

FROM alpine:latest
COPY --from=builder /myapp /myapp

COPY --from=builder /dev.env /dev.env

EXPOSE 8000

ENTRYPOINT [ "/myapp" ]