package config

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
)

// Config representa a configuração da aplicação
type Config struct {
	URL         string        `mapstructure:"url"`
	Requests    int           `mapstructure:"requests"`
	Concurrency int           `mapstructure:"concurrency"`
	Timeout     time.Duration `mapstructure:"timeout"`
}

// LoadConfig carrega a configuração a partir de flags CLI e variáveis de ambiente
func LoadConfig() (*Config, error) {
	// Definir flags
	var (
		url         = flag.String("url", "", "URL do serviço a ser testado (obrigatório)")
		requests    = flag.Int("requests", 0, "Número total de requests (obrigatório)")
		concurrency = flag.Int("concurrency", 0, "Número de chamadas simultâneas (obrigatório)")
		timeout     = flag.Duration("timeout", 30*time.Second, "Timeout para cada request")
		help        = flag.Bool("help", false, "Mostrar ajuda")
	)

	// Parsear flags
	flag.Parse()

	// Verificar se é uma solicitação de ajuda
	if *help {
		showHelp()
		os.Exit(0)
	}

	// Verificar variáveis de ambiente como fallback
	if *url == "" {
		*url = GetEnvVar("STRESSTEST_URL", "")
	}
	if *requests == 0 {
		*requests = GetEnvInt("STRESSTEST_REQUESTS", 0)
	}
	if *concurrency == 0 {
		*concurrency = GetEnvInt("STRESSTEST_CONCURRENCY", 0)
	}
	if *timeout == 30*time.Second {
		*timeout = GetEnvDuration("STRESSTEST_TIMEOUT", 30*time.Second)
	}

	// Criar configuração
	config := &Config{
		URL:         *url,
		Requests:    *requests,
		Concurrency: *concurrency,
		Timeout:     *timeout,
	}

	// Validar configuração
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

// showHelp exibe a ajuda da aplicação
func showHelp() {
	fmt.Println("Sistema CLI para realizar testes de carga em serviços web com suporte a concorrência configurável.")
	fmt.Println("")
	fmt.Println("Uso:")
	fmt.Println("  stresstest --url=URL --requests=NUM --concurrency=NUM [opções]")
	fmt.Println("")
	fmt.Println("Parâmetros obrigatórios:")
	fmt.Println("  --url string")
	fmt.Println("        URL do serviço a ser testado")
	fmt.Println("  --requests int")
	fmt.Println("        Número total de requests")
	fmt.Println("  --concurrency int")
	fmt.Println("        Número de chamadas simultâneas")
	fmt.Println("")
	fmt.Println("Parâmetros opcionais:")
	fmt.Println("  --timeout duration")
	fmt.Println("        Timeout para cada request (padrão: 30s)")
	fmt.Println("  --help")
	fmt.Println("        Mostrar esta ajuda")
	fmt.Println("")
	fmt.Println("Exemplos:")
	fmt.Println("  stresstest --url=http://google.com --requests=1000 --concurrency=10")
	fmt.Println("  stresstest --url=https://api.exemplo.com --requests=500 --concurrency=5 --timeout=60s")
	fmt.Println("")
	fmt.Println("Variáveis de ambiente:")
	fmt.Println("  STRESSTEST_URL         URL do serviço")
	fmt.Println("  STRESSTEST_REQUESTS    Número de requests")
	fmt.Println("  STRESSTEST_CONCURRENCY Número de chamadas simultâneas")
	fmt.Println("  STRESSTEST_TIMEOUT     Timeout para cada request")
}

// validateConfig valida os parâmetros da configuração
func validateConfig(config *Config) error {
	// Validar URL
	if config.URL == "" {
		return fmt.Errorf("URL é obrigatória (use --url ou STRESSTEST_URL)")
	}

	parsedURL, err := url.Parse(config.URL)
	if err != nil {
		return fmt.Errorf("URL inválida: %w", err)
	}

	if parsedURL.Scheme == "" {
		return fmt.Errorf("URL deve incluir o protocolo (http:// ou https://)")
	}

	// Validar número de requests
	if config.Requests <= 0 {
		return fmt.Errorf("número de requests deve ser maior que 0 (use --requests ou STRESSTEST_REQUESTS)")
	}

	// Validar concorrência
	if config.Concurrency <= 0 {
		return fmt.Errorf("número de concorrência deve ser maior que 0 (use --concurrency ou STRESSTEST_CONCURRENCY)")
	}

	// Validar se concorrência não é maior que requests
	if config.Concurrency > config.Requests {
		return fmt.Errorf("concorrência (%d) não pode ser maior que o número de requests (%d)", config.Concurrency, config.Requests)
	}

	// Validar timeout
	if config.Timeout <= 0 {
		return fmt.Errorf("timeout deve ser maior que 0")
	}

	return nil
}

// GetEnvVar retorna o valor de uma variável de ambiente ou o valor padrão
func GetEnvVar(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvInt retorna o valor inteiro de uma variável de ambiente ou o valor padrão
func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetEnvDuration retorna o valor de duração de uma variável de ambiente ou o valor padrão
func GetEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
