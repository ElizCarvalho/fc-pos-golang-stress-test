package reporter

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"fc-pos-golang-stress-test/internal/models"
)

// Reporter gera relatórios de teste de carga
type Reporter struct {
	stats *models.Stats
}

// NewReporter cria uma nova instância do Reporter
func NewReporter(stats *models.Stats) *Reporter {
	return &Reporter{
		stats: stats,
	}
}

// GenerateReport gera e exibe o relatório completo
func (r *Reporter) GenerateReport() {
	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Println("📊 RELATÓRIO DE TESTE DE CARGA")
	fmt.Println("=" + strings.Repeat("=", 60))

	r.printSummary()
	r.printTimingStats()
	r.printStatusCodes()
	r.printErrors()
	r.printPerformanceMetrics()

	fmt.Println("=" + strings.Repeat("=", 60))
}

// printSummary exibe o resumo geral
func (r *Reporter) printSummary() {
	fmt.Println("\n📈 RESUMO GERAL")
	fmt.Println("-" + strings.Repeat("-", 30))

	fmt.Printf("Total de Requests:     %d\n", r.stats.TotalRequests)
	fmt.Printf("Tempo Total:           %s\n", r.formatDuration(r.stats.TotalDuration))
	fmt.Printf("Requests 200 (OK):     %d (%.2f%%)\n", r.stats.Successful200, r.stats.GetSuccessRate())
	fmt.Printf("Taxa de Sucesso:       %.2f%%\n", r.stats.GetSuccessRate())
	fmt.Printf("Taxa de Erro:          %.2f%%\n", r.stats.GetErrorRate())
}

// printTimingStats exibe estatísticas de tempo
func (r *Reporter) printTimingStats() {
	fmt.Println("\n⏱️  ESTATÍSTICAS DE TEMPO")
	fmt.Println("-" + strings.Repeat("-", 30))

	fmt.Printf("Tempo Mínimo:          %s\n", r.formatDuration(r.stats.MinDuration))
	fmt.Printf("Tempo Máximo:          %s\n", r.formatDuration(r.stats.MaxDuration))
	fmt.Printf("Tempo Médio:           %s\n", r.formatDuration(r.stats.AverageDuration))

	if r.stats.TotalRequests > 0 {
		requestsPerSecond := float64(r.stats.TotalRequests) / r.stats.TotalDuration.Seconds()
		fmt.Printf("Requests por Segundo:  %.2f req/s\n", requestsPerSecond)
	}
}

// printStatusCodes exibe distribuição de códigos de status
func (r *Reporter) printStatusCodes() {
	if len(r.stats.StatusCodes) == 0 {
		return
	}

	fmt.Println("\n📋 DISTRIBUIÇÃO DE STATUS CODES")
	fmt.Println("-" + strings.Repeat("-", 30))

	// Ordenar status codes
	var statusCodes []int
	for code := range r.stats.StatusCodes {
		statusCodes = append(statusCodes, code)
	}
	sort.Ints(statusCodes)

	for _, code := range statusCodes {
		count := r.stats.StatusCodes[code]
		percentage := float64(count) / float64(r.stats.TotalRequests) * 100
		statusText := r.getStatusText(code)
		fmt.Printf("Status %d (%s):     %d (%.2f%%)\n", code, statusText, count, percentage)
	}
}

// printErrors exibe erros encontrados
func (r *Reporter) printErrors() {
	if len(r.stats.Errors) == 0 {
		return
	}

	fmt.Println("\n❌ ERROS ENCONTRADOS")
	fmt.Println("-" + strings.Repeat("-", 30))

	// Ordenar erros por frequência
	type errorCount struct {
		error string
		count int
	}

	var errors []errorCount
	for err, count := range r.stats.Errors {
		errors = append(errors, errorCount{err, count})
	}

	sort.Slice(errors, func(i, j int) bool {
		return errors[i].count > errors[j].count
	})

	for _, err := range errors {
		percentage := float64(err.count) / float64(r.stats.TotalRequests) * 100
		fmt.Printf("%s: %d (%.2f%%)\n", err.error, err.count, percentage)
	}
}

// printPerformanceMetrics exibe métricas de performance
func (r *Reporter) printPerformanceMetrics() {
	if r.stats.TotalRequests == 0 {
		return
	}

	fmt.Println("\n🚀 MÉTRICAS DE PERFORMANCE")
	fmt.Println("-" + strings.Repeat("-", 30))

	// Calcular percentis (aproximados)
	percentiles := r.calculatePercentiles()

	fmt.Printf("P50 (Mediana):         %s\n", r.formatDuration(percentiles[50]))
	fmt.Printf("P90:                   %s\n", r.formatDuration(percentiles[90]))
	fmt.Printf("P95:                   %s\n", r.formatDuration(percentiles[95]))
	fmt.Printf("P99:                   %s\n", r.formatDuration(percentiles[99]))

	// Throughput
	throughput := float64(r.stats.TotalRequests) / r.stats.TotalDuration.Seconds()
	fmt.Printf("Throughput:            %.2f req/s\n", throughput)
}

// calculatePercentiles calcula percentis aproximados
func (r *Reporter) calculatePercentiles() map[int]time.Duration {
	percentiles := make(map[int]time.Duration)

	if len(r.stats.Results) == 0 {
		return percentiles
	}

	// Extrair durações e ordenar
	var durations []time.Duration
	for _, result := range r.stats.Results {
		if result.Error == nil { // Apenas requests bem-sucedidos
			durations = append(durations, result.Duration)
		}
	}

	if len(durations) == 0 {
		return percentiles
	}

	sort.Slice(durations, func(i, j int) bool {
		return durations[i] < durations[j]
	})

	// Calcular percentis
	percentileValues := []int{50, 90, 95, 99}
	for _, p := range percentileValues {
		index := int(float64(len(durations)-1) * float64(p) / 100.0)
		if index >= 0 && index < len(durations) {
			percentiles[p] = durations[index]
		}
	}

	return percentiles
}

// getStatusText retorna uma descrição amigável do status code
func (r *Reporter) getStatusText(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "Sucesso"
	case code >= 300 && code < 400:
		return "Redirecionamento"
	case code >= 400 && code < 500:
		return "Erro do Cliente"
	case code >= 500 && code < 600:
		return "Erro do Servidor"
	default:
		return "Desconhecido"
	}
}

// formatDuration formata duração de forma legível
func (r *Reporter) formatDuration(d time.Duration) string {
	if d < time.Microsecond {
		return fmt.Sprintf("%.0fns", float64(d.Nanoseconds()))
	} else if d < time.Millisecond {
		return fmt.Sprintf("%.2fμs", float64(d.Nanoseconds())/1000)
	} else if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d.Nanoseconds())/1000000)
	} else {
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
}
