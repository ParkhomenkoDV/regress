package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"regress/internal/config"
	"regress/internal/handlers/compare"
	"regress/internal/handlers/export"
	"regress/internal/handlers/filter"
	"regress/internal/handlers/read"
	"regress/internal/shared"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-sigChan:
			cancel()
		case <-ctx.Done():
		}
	}()

	fmt.Printf("[%v] Чтение конфигурации...\n", time.Now().Format("2006-01-02 15:04:05"))

	cfg, err := config.New()
	if err != nil {
		fmt.Printf("Ошибка конфигурации %+v: %v", cfg, err)
		os.Exit(1)
	}

	fmt.Printf("[%v] Поиск файлов для сравнения...\n", time.Now().Format("2006-01-02 15:04:05"))

	befores, err := read.JSONFiles(cfg.Before)
	if err != nil {
		fmt.Printf("Ошибка поиска файлов в %s: %v", cfg.Before, err)
		os.Exit(1)
	}
	afters, err := read.JSONFiles(cfg.After)
	if err != nil {
		fmt.Printf("Ошибка поиска файлов в %s: %v", cfg.After, err)
		os.Exit(1)
	}

	fmt.Printf("[%v] Регресс...\n", time.Now().Format("2006-01-02 15:04:05"))

	// Запускаем обработку
	comparisons := compare.JSONs(ctx, befores, afters, cfg.Workers)

	fmt.Printf("[%v] Файлов с изменениями/всего: %d/%d\n", time.Now().Format("2006-01-02 15:04:05"), countChanged(comparisons), len(comparisons))

	filtered := filter.ShowAll(comparisons, cfg.ShowAll) // Фильтруем, если нужно показывать только изменения

	fmt.Printf("[%v] Экспорт в excel...\n", time.Now().Format("2006-01-02 15:04:05"))

	err = export.Excel(filtered, "comparison.xlsx")
	if err != nil {
		fmt.Printf("ошибка экспорта в excel: %v", err)
		return
	}

	fmt.Printf("[%v] Экспорт в html...\n", time.Now().Format("2006-01-02 15:04:05"))

	err = export.HTML(filtered, "comparison.html")
	if err != nil {
		fmt.Printf("ошибка экспорта в html: %v", err)
		return
	}

	fmt.Printf("[%v] Готово!\n", time.Now().Format("2006-01-02 15:04:05"))
}

func countChanged(comparisons []shared.Comparison) int {
	count := 0
	for _, comparison := range comparisons {
		if comparison.ExistsInBoth() && len(comparison.Differences) > 0 {
			count++
		}
	}
	return count
}
