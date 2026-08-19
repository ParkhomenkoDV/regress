package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path"
	"regress/internal/config"
	"regress/internal/handlers/compare"
	"regress/internal/handlers/export"
	"regress/internal/handlers/filter"
	"regress/internal/handlers/read"
	"regress/internal/shared"
	"strings"
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

	fmt.Printf("[%v] Чтение конфигурации...\n", time.Now().Format(time.DateTime))

	cfg, err := config.New()
	if err != nil {
		fmt.Printf("Ошибка конфигурации %+v: %v", cfg, err)
		os.Exit(1)
	}

	fmt.Printf("[%v] Поиск файлов для сравнения...\n", time.Now().Format(time.DateTime))

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

	fmt.Printf("[%v] Регресс...\n", time.Now().Format(time.DateTime))

	// Запускаем обработку
	comparisons := compare.JSONs(ctx, befores, afters, cfg.Workers)

	fmt.Printf("[%v] Файлов с изменениями/всего: %d/%d\n", time.Now().Format(time.DateTime), countChanged(comparisons), len(comparisons))

	fmt.Printf("[%v] Экспорт в html...\n", time.Now().Format(time.DateTime))

	folder := fmt.Sprintf("comparison_%s_%s_%s", path.Base(cfg.Before), path.Base(cfg.After), strings.ReplaceAll(time.Now().Format(time.DateTime), ":", "-"))

	err = export.HTML(comparisons, folder)
	if err != nil {
		fmt.Printf("ошибка экспорта в html: %v", err)
		return
	}

	if !cfg.ShowAll { // Фильтруем, если нужно показывать только изменения
		comparisons = filter.Diff(comparisons)
	}

	fmt.Printf("[%v] Экспорт в excel...\n", time.Now().Format(time.DateTime))

	err = export.Excel(comparisons, path.Join(folder, "comparison.xlsx"))
	if err != nil {
		fmt.Printf("ошибка экспорта в excel: %v", err)
		return
	}

	fmt.Printf("[%v] Готово!\n", time.Now().Format(time.DateTime))
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
