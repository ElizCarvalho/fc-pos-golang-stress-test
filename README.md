# 🚀 FC - POS Golang Stress Test

> Sistema CLI em Go para realizar testes de carga em serviços web com suporte a concorrência configurável.

## 📌 Sobre

Ferramenta de linha de comando desenvolvida em Go para executar testes de carga em serviços web. Permite configurar o número total de requests, o nível de concorrência e gera relatórios detalhados sobre a performance.

### Características

- ✅ **Interface CLI intuitiva** com flags configuráveis
- ✅ **Controle de concorrência** com worker pool pattern
- ✅ **Relatórios detalhados** com métricas de performance
- ✅ **Suporte a Docker** para execução em containers
- ✅ **Cancelamento graceful** com Ctrl+C

## 🔧 Instalação

```bash
# Clone e configure
git clone https://github.com/ElizCarvalho/fc-pos-golang-stress-test.git
cd fc-pos-golang-stress-test
make setup
make build
```

## 📚 Uso

### Execução Local

```bash
# Uso básico
./stresstest --url=http://google.com --requests=1000 --concurrency=10

# Com timeout customizado
./stresstest --url=https://api.exemplo.com --requests=500 --concurrency=5 --timeout=60s

# Usando variáveis de ambiente
export STRESSTEST_URL=http://google.com
export STRESSTEST_REQUESTS=1000
export STRESSTEST_CONCURRENCY=10
./stresstest
```

### Execução com Docker

```bash
# Build e execução
make docker-build
make docker-run-example

# Ou diretamente
docker run --rm fc-pos-golang-stress-test:latest --url=http://google.com --requests=1000 --concurrency=10
```

### Parâmetros

| Parâmetro | Flag | Descrição | Obrigatório | Padrão |
|-----------|------|-----------|-------------|---------|
| `--url` | `-u` | URL do serviço a ser testado | ✅ | - |
| `--requests` | `-r` | Número total de requests | ✅ | - |
| `--concurrency` | `-c` | Número de chamadas simultâneas | ✅ | - |
| `--timeout` | `-t` | Timeout para cada request | ❌ | 30s |

## 🧪 Exemplos

```bash
# Teste básico
./stresstest --url=http://google.com --requests=100 --concurrency=5

# Teste de API REST
./stresstest --url=https://jsonplaceholder.typicode.com/posts --requests=200 --concurrency=10

# Demonstração rápida
make demo
```

## 📊 Relatório de Saída

O sistema gera um relatório detalhado contendo:

- **Resumo Geral**: Total de requests, tempo total, taxa de sucesso/erro
- **Estatísticas de Tempo**: Min, max, médio, throughput
- **Distribuição de Status Codes**: Contagem e percentual de cada status HTTP
- **Erros Encontrados**: Lista de erros com frequência
- **Métricas de Performance**: Percentis (P50, P90, P95, P99)

### Exemplo de Saída

```bash
============================================================
📊 RELATÓRIO DE TESTE DE CARGA
============================================================

📈 RESUMO GERAL
------------------------------
Total de Requests:     1000
Tempo Total:           2.45s
Requests 200 (OK):     987 (98.70%)
Taxa de Sucesso:       98.70%
Taxa de Erro:          1.30%

⏱️  ESTATÍSTICAS DE TEMPO
------------------------------
Tempo Mínimo:          45.23ms
Tempo Máximo:          1.23s
Tempo Médio:           156.78ms
Requests por Segundo:  408.16 req/s

📋 DISTRIBUIÇÃO DE STATUS CODES
------------------------------
Status 200 (Sucesso):     987 (98.70%)
Status 301 (Redirecionamento): 10 (1.00%)
Status 404 (Erro do Cliente): 3 (0.30%)

❌ ERROS ENCONTRADOS
------------------------------
context deadline exceeded: 10 (1.00%)
connection refused: 3 (0.30%)

🚀 MÉTRICAS DE PERFORMANCE
------------------------------
P50 (Mediana):         142.15ms
P90:                   245.67ms
P95:                   312.45ms
P99:                   567.89ms
Throughput:            408.16 req/s
============================================================
```

## 🐳 Docker

```bash
# Build e execução
make docker-build
make docker-run

# Ou diretamente
docker run --rm fc-pos-golang-stress-test:latest --url=http://google.com --requests=100 --concurrency=10
```

## 🔧 Comandos Úteis

```bash
make setup    # Configura ambiente
make build    # Compila binário
make demo     # Demonstração rápida
make clean    # Limpa arquivos
```

---

**Desenvolvido para o desafio técnico da FC POS Golang** 🚀
