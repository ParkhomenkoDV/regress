package read

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestJSONFiles(t *testing.T) {
	// Создаём временную директорию
	tmpDir := t.TempDir()

	// Подготовка тестовых файлов и поддиректорий
	testFiles := []string{
		"file1.json",
		"file2.JSON",
		"file3.JsOn",
		"file4.txt",
		"file5",
		".hidden.json", // файл с именем .hidden.json
		"file6.json.bak",
	}
	for _, f := range testFiles {
		path := filepath.Join(tmpDir, f)
		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Поддиректория с .json файлом – должна быть проигнорирована
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "sub.json"), []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	// Ожидаемый результат: только файлы с расширением .json (регистронезависимо)
	expected := map[string]string{
		"file1.json":   filepath.Join(tmpDir, "file1.json"),
		"file2.JSON":   filepath.Join(tmpDir, "file2.JSON"),
		"file3.JsOn":   filepath.Join(tmpDir, "file3.JsOn"),
		".hidden.json": filepath.Join(tmpDir, ".hidden.json"),
		// "file6.json.bak" – не .json, не включаем
		// файлы из поддиректории не включаем
	}

	got, err := JSONFiles(tmpDir)
	if err != nil {
		t.Fatalf("не ожидалась ошибка, получена: %v", err)
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("JSONFiles() вернул %v, ожидалось %v", got, expected)
	}

	// Тест: несуществующая директория
	_, err = JSONFiles(filepath.Join(tmpDir, "nonexistent"))
	if err == nil {
		t.Error("ожидалась ошибка для несуществующей директории, получено nil")
	}

	// Тест: пустая директория
	emptyDir := t.TempDir()
	gotEmpty, err := JSONFiles(emptyDir)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if len(gotEmpty) != 0 {
		t.Errorf("для пустой директории ожидалась пустая карта, получено %v", gotEmpty)
	}
}
