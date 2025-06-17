FROM golang:1.24-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o terraform-ansible-inventory

FROM alpine:3.19
COPY --from=builder /src/terraform-ansible-inventory /usr/local/bin/terraform-ansible-inventory
ENTRYPOINT ["terraform-ansible-inventory"]
