package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"fc-pos-golang-stress-test/internal/config"
	"fc-pos-golang-stress-test/internal/executor"
	"fc-pos-golang-stress-test/internal/reporter"
)

func main() {
	// Configurar context com cancelamento
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Configurar captura de sinais para cancelamento graceful
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n🛑 Cancelamento solicitado. Finalizando...")
		cancel()
	}()

	// Carregar configuração
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Erro na configuração: %v\n", err)
		os.Exit(1)
	}

	// Criar executor
	exec := executor.NewExecutor(cfg)

	// Executar teste de carga
	fmt.Println("🔧 Iniciando sistema de teste de carga...")
	stats, err := exec.Execute(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Erro durante execução: %v\n", err)
		os.Exit(1)
	}

	// Verificar se foi cancelado
	if ctx.Err() != nil {
		fmt.Println("⚠️  Teste cancelado pelo usuário")
		os.Exit(0)
	}

	// Gerar relatório
	rep := reporter.NewReporter(stats)
	rep.GenerateReport()

	// Determinar código de saída baseado no sucesso
	if stats.GetSuccessRate() < 50.0 {
		fmt.Println("\n⚠️  Taxa de sucesso baixa detectada!")
		os.Exit(1)
	}

	fmt.Println("\n✅ Teste de carga concluído com sucesso!")
}
