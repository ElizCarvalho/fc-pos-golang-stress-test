package executor

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"fc-pos-golang-stress-test/internal/config"
	"fc-pos-golang-stress-test/internal/models"
)

// Executor gerencia a execução dos testes de carga
type Executor struct {
	config *config.Config
	client *http.Client
	stats  *models.Stats
}

// NewExecutor cria uma nova instância do Executor
func NewExecutor(cfg *config.Config) *Executor {
	client := &http.Client{
		Timeout: cfg.Timeout,
		Transport: &http.Transport{
			MaxIdleConns:        cfg.Concurrency,
			MaxIdleConnsPerHost: cfg.Concurrency,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	return &Executor{
		config: cfg,
		client: client,
		stats:  models.NewStats(),
	}
}

// Execute executa o teste de carga
func (e *Executor) Execute(ctx context.Context) (*models.Stats, error) {
	fmt.Printf("🚀 Iniciando teste de carga...\n")
	fmt.Printf("   URL: %s\n", e.config.URL)
	fmt.Printf("   Requests: %d\n", e.config.Requests)
	fmt.Printf("   Concorrência: %d\n", e.config.Concurrency)
	fmt.Printf("   Timeout: %s\n", e.config.Timeout)
	fmt.Println()

	startTime := time.Now()

	// Canal para distribuir requests entre workers
	requestChan := make(chan int, e.config.Requests)
	resultChan := make(chan models.RequestResult, e.config.Requests)

	// Preencher canal com requests
	for i := 0; i < e.config.Requests; i++ {
		requestChan <- i
	}
	close(requestChan)

	// WaitGroup para aguardar todos os workers
	var wg sync.WaitGroup

	// Iniciar workers
	for i := 0; i < e.config.Concurrency; i++ {
		wg.Add(1)
		go e.worker(ctx, i+1, requestChan, resultChan, &wg)
	}

	// Goroutine para coletar resultados
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Coletar resultados
	for result := range resultChan {
		e.stats.AddResult(result)
	}

	// Calcular estatísticas finais
	e.stats.CalculateAverage()
	e.stats.TotalDuration = time.Since(startTime)

	return e.stats, nil
}

// worker executa requests em um worker individual
func (e *Executor) worker(ctx context.Context, workerID int, requestChan <-chan int, resultChan chan<- models.RequestResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for requestID := range requestChan {
		select {
		case <-ctx.Done():
			// Context cancelado, parar execução
			return
		default:
			result := e.makeRequest(ctx, requestID, workerID)
			resultChan <- result
		}
	}
}

// makeRequest executa uma requisição HTTP individual
func (e *Executor) makeRequest(ctx context.Context, requestID, workerID int) models.RequestResult {
	startTime := time.Now()

	// Criar request com context
	req, err := http.NewRequestWithContext(ctx, "GET", e.config.URL, nil)
	if err != nil {
		return models.RequestResult{
			StatusCode: 0,
			Duration:   time.Since(startTime),
			Error:      fmt.Errorf("erro ao criar request: %w", err),
			Timestamp:  startTime,
		}
	}

	// Adicionar headers básicos
	req.Header.Set("User-Agent", "fc-pos-golang-stress-test/1.0")
	req.Header.Set("Accept", "*/*")

	// Executar request
	resp, err := e.client.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		return models.RequestResult{
			StatusCode: 0,
			Duration:   duration,
			Error:      err,
			Timestamp:  startTime,
		}
	}

	// Fechar response body
	defer resp.Body.Close()

	return models.RequestResult{
		StatusCode: resp.StatusCode,
		Duration:   duration,
		Error:      nil,
		Timestamp:  startTime,
	}
}

// GetStats retorna as estatísticas atuais
func (e *Executor) GetStats() *models.Stats {
	return e.stats
}

// ResetStats limpa as estatísticas
func (e *Executor) ResetStats() {
	e.stats = models.NewStats()
}
