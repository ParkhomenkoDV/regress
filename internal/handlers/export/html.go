package export

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"regress/internal/shared"

	"github.com/ParkhomenkoDV/progress"
)

// Missing описывает отсутствующий файл
type Missing struct {
	FileName string
	Side     string // "before" или "after"
}

// FieldSummary используется для строки таблицы на главной странице
type FieldSummary struct {
	Field      string
	Action     action
	Changes    uint64
	DetailLink string
}

// FieldFileInfo описывает изменение одного поля в одном файле
type FieldFileInfo struct {
	FileName string
	Action   action // "change", "add", "del"
	Before   any
	After    any
}

// FieldDetail используется для страницы конкретного поля
type FieldDetail struct {
	Field   string
	Records []FieldFileInfo
}

// Report — данные для главной страницы
type Report struct {
	TotalFiles   uint
	ChangedFiles uint
	Missings     []Missing
	Fields       []FieldSummary
}

// HTML создаёт папку с HTML-отчётами о регрессе
func HTML(comparisons []shared.Comparison, folderName string) error {
	// Создаём папку для регресса
	if err := os.MkdirAll(folderName, 0755); err != nil {
		return fmt.Errorf("создание папки: %w", err)
	}
	// Создаём подпапку для детальных страниц
	detailsDir := filepath.Join(folderName, "details")
	if err := os.MkdirAll(detailsDir, 0755); err != nil {
		return fmt.Errorf("создание папки details: %w", err)
	}

	// Собираем данные по полям
	fieldMap := make(map[string][]FieldFileInfo) // field -> список изменений
	var missing []Missing

	for _, comparison := range comparisons {
		if !comparison.ExistsInBoth() {
			// Файл отсутствует в одной из сторон – добавляем в missing (без информации о полях)
			side := "unknown"
			if comparison.ExistsBefore && !comparison.ExistsAfter {
				side = before
			} else if !comparison.ExistsBefore && comparison.ExistsAfter {
				side = after
			}
			missing = append(missing, Missing{FileName: comparison.FileName, Side: side})
			continue
		}
		if len(comparison.Differences) == 0 {
			continue // нет изменений
		}
		for _, diff := range comparison.Differences {
			action := determineAction(diff.Before, diff.After)
			fieldMap[diff.Field] = append(fieldMap[diff.Field], FieldFileInfo{
				FileName: comparison.FileName,
				Action:   action,
				Before:   diff.Before,
				After:    diff.After,
			})
		}
	}

	bar := progress.New(time.Second, "🧾 HTML:", 50, uint64(len(comparisons)), true, true, false)
	cancelBar := bar.Start(context.Background())

	// Преобразуем в срез для главной страницы и генерируем файлы полей
	fields := make([]FieldSummary, 0, len(fieldMap))
	for field, records := range fieldMap {
		// Определяем общий action для поля
		action := aggregateAction(records)
		// Сортируем записи по имени файла для стабильности
		sort.Slice(records, func(i, j int) bool {
			return records[i].FileName < records[j].FileName
		})
		// Сохраняем детали в отдельный HTML-файл внутри подпапки details
		detailFileName := sanitizeFileName(field) + ".html"
		detailPath := filepath.Join(detailsDir, detailFileName)
		if err := writeFieldPage(detailPath, field, records); err != nil {
			return fmt.Errorf("ошибка создания страницы для поля %s: %w", field, err)
		}
		// Добавляем запись для главной страницы, ссылка указывает на подпапку details
		fields = append(fields, FieldSummary{
			Field:      field,
			Action:     action,
			Changes:    uint64(len(records)),
			DetailLink: "details/" + detailFileName, // относительный путь от regress.html
		})
		bar.Add(1)
	}
	cancelBar()

	// Сортируем поля по имени
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Field < fields[j].Field
	})

	// Считаем общую статистику
	var changedFiles uint
	for _, comparison := range comparisons {
		if comparison.ExistsInBoth() && len(comparison.Differences) > 0 {
			changedFiles++
		}
	}

	// Генерируем главную страницу
	data := Report{
		TotalFiles:   uint(len(comparisons)),
		ChangedFiles: changedFiles,
		Missings:     missing,
		Fields:       fields,
	}
	mainPath := filepath.Join(folderName, "regress.html")
	if err := writeMainPage(mainPath, data); err != nil {
		return fmt.Errorf("ошибка создания главной страницы: %w", err)
	}

	return nil
}

// determineAction определяет тип изменения по значениям до и после
func determineAction(before, after any) action {
	switch {
	case before == nil && after != nil:
		return add
	case before != nil && after == nil:
		return del
	default:
		return change
	}
}

// aggregateAction определяет общее действие для поля (приоритет: del > add > change)
func aggregateAction(records []FieldFileInfo) action {
	hasDel, hasAdd, hasChange := false, false, false
	for _, r := range records {
		switch r.Action {
		case "del":
			hasDel = true
		case "add":
			hasAdd = true
		default:
			hasChange = true
		}
	}
	if hasDel {
		return del
	}
	if hasAdd {
		return add
	}
	if hasChange {
		return change
	}
	return ""
}

// sanitizeFileName заменяет небезопасные символы на '_'
func sanitizeFileName(name string) string {
	// Заменяем все символы, кроме букв, цифр, точки, дефиса и подчеркивания
	return strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '.' || r == '-' || r == '_' {
			return r
		}
		return '_'
	}, name)
}

const mainTemplate = `
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <title>Regress – main</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        h1, h2 { color: #333; }
        table { border-collapse: collapse; width: 100%; margin-top: 20px; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
        .missing { color: #cc0000; }
        .badge { background: #f0f0f0; padding: 2px 8px; border-radius: 12px; font-size: 0.9em; }
        .action-change { color: #0066cc; }
        .action-add { color: #008000; }
        .action-del { color: #cc0000; }
        .stat-number { color: #0066cc; }
    </style>
</head>
<body>
    <h1>Regress</h1>
    <div class="stat">
        <p>Total: <strong>{{.TotalFiles}}</strong></p>
        <p class="stat-number">Changes: <strong>{{.ChangedFiles}}</strong></p>
        {{if .Missings}}
        <p class="missing">NMiss: <strong>{{len .Missings}}</strong></p>
        {{end}}
    </div>

    <h2>Differences</h2>
    <table>
        <thead>
            <tr><th>Action</th><th>Field</th><th>NDiff</th><th>Details</th></tr>
        </thead>
        <tbody>
        {{range .Fields}}
        <tr>
            <td><span class="action-{{.Action}}">{{.Action}}</span></td>
            <td>{{.Field}}</td>
            <td>{{.Changes}}</td>
            <td><a href="{{.DetailLink}}">more</a></td>
        </tr>
        {{else}}
        <tr><td colspan="4">Equal</td></tr>
        {{end}}
        </tbody>
    </table>

    {{if .Missings}}
    <h2>Missings</h2>
    <ul>
        {{range .Missings}}
        <li class="missing">{{.FileName}} ({{.Side}})</li>
        {{end}}
    </ul>
    {{end}}
</body>
</html>`

// writeMainPage записывает главную страницу
func writeMainPage(path string, data Report) error {
	t, err := template.New("main").Parse(mainTemplate)
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return t.Execute(file, data)
}

const fieldTemplate = `
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <title>Field: {{.Field}}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        h1 { color: #333; }
        table { border-collapse: collapse; width: 100%; margin-top: 20px; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
        .action-change { color: #0066cc; }
        .action-add { color: #008000; }
        .action-del { color: #cc0000; }
        .back { margin-bottom: 20px; }
    </style>
</head>
<body>
    <div class="back"><a href="../regress.html">← Back to regress</a></div>
    <h1>Field: {{.Field}}</h1>
    <table>
        <thead>
            <tr><th>File</th><th>Before</th><th>After</th><th>Action</th></tr>
        </thead>
        <tbody>
        {{range .Records}}
        <tr>
            <td>{{.FileName}}</td>
            <td>{{.Before}}</td>
            <td>{{.After}}</td>
            <td><span class="action-{{.Action}}">{{.Action}}</span></td>
        </tr>
        {{end}}
        </tbody>
    </table>
</body>
</html>`

// writeFieldPage записывает страницу для одного поля
func writeFieldPage(path, field string, records []FieldFileInfo) error {
	data := FieldDetail{Field: field, Records: records}
	t, err := template.New("field").Parse(fieldTemplate)
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return t.Execute(file, data)
}
