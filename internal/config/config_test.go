package config

import (
	"flag"
	"os"
	"runtime"
	"testing"
)

// TestParseFlags проверяет корректность парсинга флагов (без валидации)
func TestParseFlags(t *testing.T) {
	numCPU := runtime.NumCPU()

	tests := []struct {
		name        string
		args        []string
		wantBefore  string
		wantAfter   string
		wantShowAll bool
		wantWorkers int
	}{
		{
			name:        "все флаги заданы",
			args:        []string{"cmd", "-before", "dir1", "-after", "dir2", "-all", "-workers", "4"},
			wantBefore:  "dir1",
			wantAfter:   "dir2",
			wantShowAll: true,
			wantWorkers: 4,
		},
		{
			name:        "только обязательные флаги",
			args:        []string{"cmd", "-before", "before_dir", "-after", "after_dir"},
			wantBefore:  "before_dir",
			wantAfter:   "after_dir",
			wantShowAll: false,
			wantWorkers: numCPU,
		},
		{
			name:        "дефолтные значения (без флагов)",
			args:        []string{"cmd"},
			wantBefore:  "before",
			wantAfter:   "after",
			wantShowAll: false,
			wantWorkers: numCPU,
		},
		{
			name:        "относительные пути",
			args:        []string{"cmd", "-before", "./in", "-after", "../out"},
			wantBefore:  "./in",
			wantAfter:   "../out",
			wantShowAll: false,
			wantWorkers: numCPU,
		},
		{
			name:        "флаг all без значения",
			args:        []string{"cmd", "-before", "b", "-after", "a", "-all"},
			wantBefore:  "b",
			wantAfter:   "a",
			wantShowAll: true,
			wantWorkers: numCPU,
		},
		{
			name:        "workers задан явно",
			args:        []string{"cmd", "-before", "b", "-after", "a", "-workers", "8"},
			wantBefore:  "b",
			wantAfter:   "a",
			wantShowAll: false,
			wantWorkers: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			os.Args = tt.args
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			flags, err := parse()
			if err != nil {
				t.Fatalf("parse() вернул ошибку: %v", err)
			}

			if flags.Before != tt.wantBefore {
				t.Errorf("Before = %q, ожидалось %q", flags.Before, tt.wantBefore)
			}
			if flags.After != tt.wantAfter {
				t.Errorf("After = %q, ожидалось %q", flags.After, tt.wantAfter)
			}
			if flags.ShowAll != tt.wantShowAll {
				t.Errorf("ShowAll = %v, ожидалось %v", flags.ShowAll, tt.wantShowAll)
			}
			if flags.Workers != tt.wantWorkers {
				t.Errorf("Workers = %d, ожидалось %d", flags.Workers, tt.wantWorkers)
			}
		})
	}
}

// TestParseFlagOrder проверяет, что порядок флагов не влияет на результат
func TestParseFlagOrder(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd", "-workers", "2", "-after", "out", "-all", "-before", "in"}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	flags, err := parse()
	if err != nil {
		t.Fatalf("parse() вернул ошибку: %v", err)
	}

	if flags.Before != "in" {
		t.Errorf("Before = %q, ожидалось %q", flags.Before, "in")
	}
	if flags.After != "out" {
		t.Errorf("After = %q, ожидалось %q", flags.After, "out")
	}
	if !flags.ShowAll {
		t.Error("ShowAll должен быть true")
	}
	if flags.Workers != 2 {
		t.Errorf("Workers = %d, ожидалось %d", flags.Workers, 2)
	}
}

// TestParseDuplicateFlags проверяет, что последнее значение флага имеет приоритет
func TestParseDuplicateFlags(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd", "-before", "first", "-after", "first_after", "-workers", "1",
		"-before", "second", "-after", "second_after", "-workers", "4", "-all", "-all=false"}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	flags, err := parse()
	if err != nil {
		t.Fatalf("parse() вернул ошибку: %v", err)
	}

	if flags.Before != "second" {
		t.Errorf("Before = %q, ожидалось %q", flags.Before, "second")
	}
	if flags.After != "second_after" {
		t.Errorf("After = %q, ожидалось %q", flags.After, "second_after")
	}
	if flags.Workers != 4 {
		t.Errorf("Workers = %d, ожидалось %d", flags.Workers, 4)
	}
	// Для булевых флагов последнее значение: "-all=false" выключает флаг
	if flags.ShowAll {
		t.Error("ShowAll должен быть false (последнее значение -all=false)")
	}
}

// TestNew_ValidationErrors проверяет ошибки валидации (пустые директории, workers<1)
func TestNew_ValidationErrors(t *testing.T) {
	// Создаём временные директории для корректных путей
	tmpBefore := t.TempDir()
	tmpAfter := t.TempDir()

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "пустой Before",
			args:    []string{"cmd", "-before", "", "-after", tmpAfter},
			wantErr: true,
		},
		{
			name:    "пустой After",
			args:    []string{"cmd", "-before", tmpBefore, "-after", ""},
			wantErr: true,
		},
		{
			name:    "пустые обе",
			args:    []string{"cmd", "-before", "", "-after", ""},
			wantErr: true,
		},
		{
			name:    "workers = 0",
			args:    []string{"cmd", "-before", tmpBefore, "-after", tmpAfter, "-workers", "0"},
			wantErr: true,
		},
		{
			name:    "workers отрицательный",
			args:    []string{"cmd", "-before", tmpBefore, "-after", tmpAfter, "-workers", "-1"},
			wantErr: true,
		},
		{
			name:    "корректные значения (ошибки не должно быть)",
			args:    []string{"cmd", "-before", tmpBefore, "-after", tmpAfter, "-workers", "2"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			os.Args = tt.args
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			cfg, err := New()
			if tt.wantErr {
				if err == nil {
					t.Error("ожидалась ошибка, но не получена")
				}
				// при ошибке возвращается пустая структура
				if cfg == nil || cfg.Before != "" || cfg.After != "" {
					t.Errorf("при ошибке ожидалась пустая структура, получено %+v", cfg)
				}
			} else {
				if err != nil {
					t.Fatalf("не ожидалась ошибка, получена: %v", err)
				}
				// проверяем, что значения переданы правильно
				if cfg.Before != tmpBefore || cfg.After != tmpAfter || cfg.Workers != 2 {
					t.Errorf("некорректные значения: %+v", cfg)
				}
			}
		})
	}
}

// TestNew_ValidConfig тестирует успешное создание конфигурации с корректными значениями
func TestNew_ValidConfig(t *testing.T) {
	tmpBefore := t.TempDir()
	tmpAfter := t.TempDir()

	tests := []struct {
		name        string
		args        []string
		wantBefore  string
		wantAfter   string
		wantShowAll bool
		wantWorkers int
	}{
		{
			name:        "все флаги заданы",
			args:        []string{"cmd", "-before", tmpBefore, "-after", tmpAfter, "-all", "-workers", "4"},
			wantBefore:  tmpBefore,
			wantAfter:   tmpAfter,
			wantShowAll: true,
			wantWorkers: 4,
		},
		{
			name:        "только обязательные",
			args:        []string{"cmd", "-before", tmpBefore, "-after", tmpAfter},
			wantBefore:  tmpBefore,
			wantAfter:   tmpAfter,
			wantShowAll: false,
			wantWorkers: runtime.NumCPU(),
		},
		{
			name:        "с флагом all",
			args:        []string{"cmd", "-before", tmpBefore, "-after", tmpAfter, "-all"},
			wantBefore:  tmpBefore,
			wantAfter:   tmpAfter,
			wantShowAll: true,
			wantWorkers: runtime.NumCPU(),
		},
		{
			name:        "явно указан workers",
			args:        []string{"cmd", "-before", tmpBefore, "-after", tmpAfter, "-workers", "10"},
			wantBefore:  tmpBefore,
			wantAfter:   tmpAfter,
			wantShowAll: false,
			wantWorkers: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			os.Args = tt.args
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			cfg, err := New()
			if err != nil {
				t.Fatalf("не ожидалась ошибка, получена: %v", err)
			}

			if cfg.Before != tt.wantBefore {
				t.Errorf("Before = %q, ожидалось %q", cfg.Before, tt.wantBefore)
			}
			if cfg.After != tt.wantAfter {
				t.Errorf("After = %q, ожидалось %q", cfg.After, tt.wantAfter)
			}
			if cfg.ShowAll != tt.wantShowAll {
				t.Errorf("ShowAll = %v, ожидалось %v", cfg.ShowAll, tt.wantShowAll)
			}
			if cfg.Workers != tt.wantWorkers {
				t.Errorf("Workers = %d, ожидалось %d", cfg.Workers, tt.wantWorkers)
			}
		})
	}
}

// TestNew_WorkersDefaultCPU проверяет, что workers по умолчанию равен runtime.NumCPU()
func TestNew_WorkersDefaultCPU(t *testing.T) {
	tmpBefore := t.TempDir()
	tmpAfter := t.TempDir()

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd", "-before", tmpBefore, "-after", tmpAfter}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	cfg, err := New()
	if err != nil {
		t.Fatalf("не ожидалась ошибка: %v", err)
	}

	expected := runtime.NumCPU()
	if cfg.Workers != expected {
		t.Errorf("Workers = %d, ожидалось %d", cfg.Workers, expected)
	}
}
