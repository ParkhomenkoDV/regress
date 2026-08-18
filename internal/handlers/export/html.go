package export

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"

	"regress/internal/shared"
)

const tmpl = `
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <title>Regress Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        h1, h2 { color: #333; }
        .stat { margin-bottom: 20px; }
        .field-list { list-style: none; padding: 0; }
        .field-list li { margin: 5px 0; }
        .field-list a { text-decoration: none; color: #0066cc; cursor: pointer; }
        .field-list a:hover { text-decoration: underline; }
        .detail { display: none; margin-top: 20px; border-top: 1px solid #ccc; padding-top: 10px; }
        .detail.active { display: block; }
        table { border-collapse: collapse; width: 100%; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
        .missing { color: #cc0000; }
        .badge { background: #f0f0f0; padding: 2px 8px; border-radius: 12px; font-size: 0.9em; }
    </style>
    <script>
        function showDetail(fieldId) {
            // Скрыть все активные детали
            document.querySelectorAll('.detail').forEach(el => el.classList.remove('active'));
            // Показать нужную
            const el = document.getElementById('detail-' + fieldId);
            if (el) el.classList.add('active');
            // Обновить URL с якорем для возможности прямой ссылки
            if (history.pushState) {
                history.pushState(null, null, '#' + fieldId);
            }
        }
        // При загрузке проверяем якорь
        window.onload = function() {
            const hash = window.location.hash.substring(1);
            if (hash) {
                showDetail(hash);
            }
        };
    </script>
</head>
<body>
    <h1>Отчёт о регрессе</h1>
    <div class="stat">
        <p>Всего файлов: <strong>{{.TotalFiles}}</strong></p>
        <p>Файлов с изменениями: <strong>{{.ChangedFiles}}</strong></p>
        {{if .MissingFiles}}
        <p class="missing">Отсутствующих файлов: <strong>{{len .MissingFiles}}</strong></p>
        {{end}}
    </div>

    <h2>Список полей с изменениями</h2>
    <ul class="field-list">
        {{range .Fields}}
        <li><a onclick="showDetail('{{.Name}}'); return false;" href="#{{.Name}}">{{.Name}} <span class="badge">{{.FileCount}} файлов</span></a></li>
        {{else}}
        <li>Нет изменений.</li>
        {{end}}
    </ul>

    {{range .Fields}}
    <div id="detail-{{.Name}}" class="detail">
        <h3>Поле: {{.Name}}</h3>
        <table>
            <thead>
                <tr><th>Файл</th><th>Было</th><th>Стало</th></tr>
            </thead>
            <tbody>
                {{range .Files}}
                <tr><td>{{.FileName}}</td><td>{{.Before}}</td><td>{{.After}}</td></tr>
                {{end}}
            </tbody>
        </table>
    </div>
    {{end}}

    {{if .MissingFiles}}
    <h2>Отсутствующие файлы</h2>
    <ul>
        {{range .MissingFiles}}
        <li class="missing">{{.FileName}} ({{.Side}})</li>
        {{end}}
    </ul>
    {{end}}
</body>
</html>
`

// FieldChange описывает изменение одного поля в одном файле
type FieldChange struct {
	FileName string
	Before   string
	After    string
}

// FieldReport содержит информацию по одному полю
type FieldReport struct {
	Name      string
	FileCount int
	Files     []FieldChange
}

// MissingFile описывает отсутствующий файл
type MissingFile struct {
	FileName string
	Side     string // "before" или "after"
}

// ReportData — данные для шаблона
type ReportData struct {
	TotalFiles   int
	ChangedFiles int
	MissingFiles []MissingFile
	Fields       []FieldReport
}

// HTML создаёт HTML-отчёт о регрессе
func HTML(comparisons []shared.Comparison, fileName string) error {
	// 1. Собираем статистику и изменения по полям
	fieldMap := make(map[string]map[string]FieldChange) // field -> (fileName -> FieldChange)
	var missing []MissingFile

	for _, comp := range comparisons {
		if !comp.ExistsInBoth() {
			side := "after"
			_ = side
			// определить сторону можно по наличию файла? В текущей структуре нет информации,
			// поэтому будем считать, что если файла нет в after, то он только в before.
			// Но у нас нет явного указания, в какой из директорий файл отсутствует.
			// В структуре Comparison нет поля, указывающего, с какой стороны отсутствует.
			// Придётся сделать допущение: если ExistsInBoth == false, мы не знаем, где файл отсутствует.
			// В реальности такая ситуация возникает, если файл есть только в before или только в after.
			// Можно передать дополнительную информацию, но сейчас просто отметим как "missing".
			missing = append(missing, MissingFile{FileName: comp.FileName, Side: "unknown"})
			continue
		}
		if len(comp.Differences) == 0 {
			continue // нет изменений
		}
		// Для каждого различия добавляем в карту
		for _, diff := range comp.Differences {
			field := diff.Field
			if _, ok := fieldMap[field]; !ok {
				fieldMap[field] = make(map[string]FieldChange)
			}
			// Если для этого файла уже есть запись по этому полю (маловероятно, но перезапишем)
			fieldMap[field][comp.FileName] = FieldChange{
				FileName: comp.FileName,
				Before:   fmt.Sprint(diff.Before),
				After:    fmt.Sprint(diff.After),
			}
		}
	}

	// 2. Преобразуем карту в срез FieldReport
	fields := make([]FieldReport, 0, len(fieldMap))
	for field, fileMap := range fieldMap {
		files := make([]FieldChange, 0, len(fileMap))
		for _, fc := range fileMap {
			files = append(files, fc)
		}
		// сортируем по имени файла для стабильности
		sort.Slice(files, func(i, j int) bool {
			return files[i].FileName < files[j].FileName
		})
		fields = append(fields, FieldReport{
			Name:      field,
			FileCount: len(files),
			Files:     files,
		})
	}
	// Сортируем поля по имени
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})

	// 3. Считаем общую статистику
	totalFiles := len(comparisons)
	changedFiles := 0
	for _, comp := range comparisons {
		if comp.ExistsInBoth() && len(comp.Differences) > 0 {
			changedFiles++
		}
	}

	// 4. Подготавливаем данные для шаблона
	data := ReportData{
		TotalFiles:   totalFiles,
		ChangedFiles: changedFiles,
		MissingFiles: missing,
		Fields:       fields,
	}

	// 5. Загружаем и выполняем шаблон
	t, err := template.New("report").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("парсинг шаблона: %w", err)
	}

	// Создаём директорию, если её нет
	dir := filepath.Dir(fileName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("создание директории: %w", err)
	}

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("создание файла: %w", err)
	}
	defer file.Close()

	if err := t.Execute(file, data); err != nil {
		return fmt.Errorf("выполнение шаблона: %w", err)
	}

	return nil
}
