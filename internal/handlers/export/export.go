package export

import (
	"regress/internal/shared"

	excel "github.com/xuri/excelize/v2"
)

func Excel(comparisons []shared.Comparison, filename string) error {
	f := excel.NewFile() // Создаем новый Excel файл

	// Получаем все уникальные названия полей для создания заголовков
	fieldSet := make(map[string]bool)
	for _, comp := range comparisons {
		for _, diff := range comp.Differences {
			fieldSet[diff.Field] = true
		}
	}

	// Преобразуем в упорядоченный список полей
	fields := make([]string, 0, len(fieldSet))
	for field := range fieldSet {
		fields = append(fields, field)
	}

	// Создаем заголовки
	headers := []string{"FileName", "ExistsInBoth"}
	for _, field := range fields {
		headers = append(headers, field+"_before", field+"_after")
	}

	// Записываем заголовки в первую строку
	for col, header := range headers {
		cell, _ := excel.CoordinatesToCellName(col+1, 1)
		f.SetCellValue("Sheet1", cell, header)
	}

	// Записываем данные
	for row, comp := range comparisons {
		// Преобразуем различия в карту для быстрого доступа
		diffMap := make(map[string]shared.Difference)
		for _, diff := range comp.Differences {
			diffMap[diff.Field] = diff
		}

		// Заполняем ячейки
		col := 1

		// FileName
		cell, _ := excel.CoordinatesToCellName(col, row+2)
		f.SetCellValue("Sheet1", cell, comp.FileName)
		col++

		// ExistsInBoth
		cell, _ = excel.CoordinatesToCellName(col, row+2)
		f.SetCellValue("Sheet1", cell, comp.ExistsInBoth)
		col++

		// Данные для каждого поля
		for _, field := range fields {
			if diff, exists := diffMap[field]; exists {
				// Before значение
				cell, _ = excel.CoordinatesToCellName(col, row+2)
				f.SetCellValue("Sheet1", cell, diff.Before)
				col++

				// After значение
				cell, _ = excel.CoordinatesToCellName(col, row+2)
				f.SetCellValue("Sheet1", cell, diff.After)
				col++
			} else {
				// Если поля нет в различиях, оставляем пустые ячейки
				col += 2
			}
		}
	}
	// Настраиваем ширину колонок для лучшего отображения
	for i := 1; i <= len(headers); i++ {
		f.SetColWidth("Sheet1", string(rune('A'+i-1)), string(rune('A'+i-1)), 20)
	}

	// Записываем форматирование для заголовков
	style, _ := f.NewStyle(&excel.Style{
		Font: &excel.Font{
			Bold: true,
		},
		Fill: excel.Fill{
			Type:    "pattern",
			Color:   []string{"#E0E0E0"},
			Pattern: 1,
		},
	})

	// Применяем стиль к заголовкам
	for i := 1; i <= len(headers); i++ {
		cell, _ := excel.CoordinatesToCellName(i, 1)
		f.SetCellStyle("Sheet1", cell, cell, style)
	}

	// Сохраняем файл
	if err := f.SaveAs(filename); err != nil {
		return err
	}

	return nil
}
