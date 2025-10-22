# ==============================================================================
# Multi-stage Dockerfile para fc-pos-golang-stress-test
# ==============================================================================

# Estágio 1: Build
FROM golang:1.23-alpine AS builder

# Instalar dependências necessárias para build
RUN apk add --no-cache ca-certificates git tzdata

# Definir diretório de trabalho
WORKDIR /app

# Copiar arquivos de dependências
COPY go.mod go.sum ./

# Baixar dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Build da aplicação
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o stresstest \
    ./cmd/stresstest

# Estágio 2: Runtime
FROM alpine:3.18

# Instalar certificados CA e timezone data, criar usuário e configurar permissões
RUN apk --no-cache add ca-certificates tzdata && \
    addgroup -g 1001 -S stresstest && \
    adduser -u 1001 -S stresstest -G stresstest

# Definir diretório de trabalho
WORKDIR /app

# Copiar binário do estágio de build e definir permissões
COPY --from=builder /app/stresstest /app/stresstest
RUN chmod +x /app/stresstest

# Mudar para usuário não-root
USER stresstest

# Expor porta (opcional, para health checks)
EXPOSE 8080

# Definir entrypoint
ENTRYPOINT ["/app/stresstest"]

# Labels para metadados
LABEL maintainer="fc-pos-golang-stress-test"
LABEL version="1.0.0"
LABEL description="Sistema CLI para testes de carga em serviços web"
