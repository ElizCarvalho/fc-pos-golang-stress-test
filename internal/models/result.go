package models

import (
	"time"
)

// RequestResult representa o resultado de uma requisição individual
type RequestResult struct {
	StatusCode int           `json:"status_code"`
	Duration   time.Duration `json:"duration"`
	Error      error         `json:"error,omitempty"`
	Timestamp  time.Time     `json:"timestamp"`
}

// Stats representa as estatísticas agregadas de todos os requests
type Stats struct {
	TotalRequests   int             `json:"total_requests"`
	TotalDuration   time.Duration   `json:"total_duration"`
	Successful200   int             `json:"successful_200"`
	StatusCodes     map[int]int     `json:"status_codes"`
	Errors          map[string]int  `json:"errors"`
	MinDuration     time.Duration   `json:"min_duration"`
	MaxDuration     time.Duration   `json:"max_duration"`
	AverageDuration time.Duration   `json:"average_duration"`
	Results         []RequestResult `json:"results"`
}

// NewStats cria uma nova instância de Stats
func NewStats() *Stats {
	return &Stats{
		StatusCodes: make(map[int]int),
		Errors:      make(map[string]int),
		Results:     make([]RequestResult, 0),
	}
}

// AddResult adiciona um resultado de request às estatísticas
func (s *Stats) AddResult(result RequestResult) {
	s.Results = append(s.Results, result)
	s.TotalRequests++

	// Atualizar duração total
	s.TotalDuration += result.Duration

	// Atualizar min/max duration
	if s.MinDuration == 0 || result.Duration < s.MinDuration {
		s.MinDuration = result.Duration
	}
	if result.Duration > s.MaxDuration {
		s.MaxDuration = result.Duration
	}

	// Contar status codes
	if result.Error != nil {
		errorType := result.Error.Error()
		s.Errors[errorType]++
	} else {
		s.StatusCodes[result.StatusCode]++
		if result.StatusCode == 200 {
			s.Successful200++
		}
	}
}

// CalculateAverage calcula a duração média dos requests
func (s *Stats) CalculateAverage() {
	if s.TotalRequests > 0 {
		s.AverageDuration = s.TotalDuration / time.Duration(s.TotalRequests)
	}
}

// GetSuccessRate retorna a taxa de sucesso (requests 200 / total)
func (s *Stats) GetSuccessRate() float64 {
	if s.TotalRequests == 0 {
		return 0
	}
	return float64(s.Successful200) / float64(s.TotalRequests) * 100
}

// GetErrorRate retorna a taxa de erro (requests com erro / total)
func (s *Stats) GetErrorRate() float64 {
	if s.TotalRequests == 0 {
		return 0
	}
	totalErrors := 0
	for _, count := range s.Errors {
		totalErrors += count
	}
	return float64(totalErrors) / float64(s.TotalRequests) * 100
}
