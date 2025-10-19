# ==============================================================================
# Makefile para fc-pos-golang-stress-test
# ==============================================================================

APP_NAME=stresstest
DOCKER_IMAGE=fc-pos-golang-stress-test

.PHONY: setup build demo docker-build docker-run clean help

setup: ## Configura o ambiente
	@echo "🔧 Configurando ambiente..."
	@go mod download
	@go mod tidy
	@echo "✅ Ambiente configurado!"

build: ## Compila o binário
	@echo "🔨 Compilando binário..."
	@go build -o $(APP_NAME) ./cmd/stresstest
	@echo "✅ Binário compilado: $(APP_NAME)"

demo: build ## Demonstração básica
	@echo "🎯 Executando demonstração..."
	@./$(APP_NAME) --url=http://google.com --requests=50 --concurrency=5

docker-build: ## Build da imagem Docker
	@echo "🐳 Construindo imagem Docker..."
	@docker build -t $(DOCKER_IMAGE):latest .
	@echo "✅ Imagem construída: $(DOCKER_IMAGE):latest"

docker-run: docker-build ## Executa exemplo via Docker
	@echo "🐳 Executando exemplo via Docker..."
	@docker run --rm $(DOCKER_IMAGE):latest --url=http://google.com --requests=100 --concurrency=10

clean: ## Limpa arquivos gerados
	@echo "🧹 Limpando arquivos..."
	@rm -f $(APP_NAME)
	@go clean
	@echo "✅ Limpeza concluída!"

help: ## Mostra esta ajuda
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
