FROM golang:1.24-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o terraform-ansible-inventory

FROM alpine:3.19 AS runtime
RUN apk add --no-cache git curl vim
COPY --from=builder /src/terraform-ansible-inventory /usr/local/bin/terraform-ansible-inventory
ENTRYPOINT ["terraform-ansible-inventory"]

FROM runtime AS terraform
RUN apk add --no-cache unzip \
    && TF_VERSION=$(curl -fsSL https://checkpoint-api.hashicorp.com/v1/check/terraform | grep -o '"current_version":"[^"]*"' | cut -d '"' -f4) \
    && wget -q https://releases.hashicorp.com/terraform/${TF_VERSION}/terraform_${TF_VERSION}_linux_amd64.zip -O /tmp/terraform.zip \
    && unzip /tmp/terraform.zip -d /usr/local/bin \
    && rm /tmp/terraform.zip
