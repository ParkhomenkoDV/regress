package read

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestJSONFiles(t *testing.T) {
	// 1. Создаём временную директорию с файлами
	tmpDir := t.TempDir()

	testFiles := []string{
		"a.json",
		"b.JSON",
		"c.JsOn",
		"d.txt",
		"e",
		".json",      // файл с именем .json (расширение и имя совпадают)
		"f.json.bak", // не .json
	}
	for _, f := range testFiles {
		path := filepath.Join(tmpDir, f)
		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// 2. Создаём поддиректорию с .json файлом (должен быть проигнорирован)
	subDir := filepath.Join(tmpDir, "sub")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	subFile := filepath.Join(subDir, "sub.json")
	if err := os.WriteFile(subFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	expected := []string{"a.json", "b.JSON", "c.JsOn", ".json"}
	sort.Strings(expected) // порядок не гарантируется

	got, err := JSONFiles(tmpDir)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	sort.Strings(got)

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("ожидалось %v, получено %v", expected, got)
	}

	// 3. Тест на несуществующую директорию
	_, err = JSONFiles(filepath.Join(tmpDir, "nonexistent"))
	if err == nil {
		t.Error("ожидалась ошибка для несуществующей директории, получено nil")
	}

	// 4. Тест на пустую директорию
	emptyDir := t.TempDir()
	got, err = JSONFiles(emptyDir)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("ожидался пустой слайс, получено %v", got)
	}
}
