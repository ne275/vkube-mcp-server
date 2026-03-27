# 构建：在 vkube-mcp 目录下执行
#   docker build -t tzr0125/vkube-mcp:demo .
# Apple Silicon / ARM 本机构建 x86_64（amd64）并推送 Hub：
#   docker buildx build --platform linux/amd64 -t tzr0125/vkube-mcp:demo --push .
# Podman：
#   podman build --platform linux/amd64 -t docker.io/tzr0125/vkube-mcp:demo .
#   podman push docker.io/tzr0125/vkube-mcp:demo
# 运行示例：
#   docker run --rm -p 3100:3100 -e VKUBE_MCP_PUBLIC_BASE_URL=http://host.docker.internal:3100 tzr0125/vkube-mcp:demo

FROM golang:1.22-alpine AS builder
RUN apk add --no-cache git ca-certificates
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/vkube-mcp .

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /out/vkube-mcp /vkube-mcp

ENV VKUBE_MCP_LISTEN=:3100
EXPOSE 3100
USER nonroot:nonroot
ENTRYPOINT ["/vkube-mcp"]
