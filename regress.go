package main

import (
	"fmt"
	"os"
	"regress/internal/config"
	"regress/internal/handlers/compare"
	"regress/internal/handlers/export"
	"regress/internal/shared"
	"time"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Printf("Ошибка конфигурации %+v: %v", cfg, err)
		os.Exit(1)
	}

	fmt.Printf("[%v] Регресс запущен...\n", time.Now().Format("2006-01-02 15:04:05"))

	comparisons, err := compare.JSONs(cfg.Before, cfg.After, cfg.Workers)
	if err != nil {
		fmt.Printf("ошибка сравнения json: %v", err)
		os.Exit(1)
	}

	fmt.Printf("[%v] Файлов с изменениями/всего: %d/%d\n", time.Now().Format("2006-01-02 15:04:05"), countChanged(comparisons), len(comparisons))

	filtered := filterComparisons(comparisons, cfg.ShowAll) // Фильтруем, если нужно показывать только изменения

	fmt.Printf("[%v] Экспорт в excel...\n", time.Now().Format("2006-01-02 15:04:05"))

	err = export.Excel(filtered, "comparison.xlsx")
	if err != nil {
		fmt.Printf("ошибка экспорта в excel: %v", err)
		return
	}

	fmt.Printf("[%v] Регрес готов!\n", time.Now().Format("2006-01-02 15:04:05"))
}

func countChanged(comparisons []shared.Comparison) int {
	count := 0
	for _, comp := range comparisons {
		if comp.ExistsInBoth && len(comp.Differences) > 0 {
			count++
		}
	}
	return count
}

// filterComparisons фильтрует сравнения по настройкам
func filterComparisons(comparisons []shared.Comparison, showAll bool) []shared.Comparison {
	if showAll {
		return comparisons
	}

	filtered := make([]shared.Comparison, 0, len(comparisons))
	for _, comp := range comparisons {
		if len(comp.Differences) > 0 || !comp.ExistsInBoth {
			filtered = append(filtered, comp)
		}
	}
	return filtered
}
