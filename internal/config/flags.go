package config

import (
	"flag"
	"fmt"
	"runtime"
)

const usage = "Использование: go run regress.go [--before <путь>] [--after <путь>] [--all] [--workers N]"

type Flags struct {
	BeforeDir string
	AfterDir  string
	ShowAll   bool
	Workers   int
}

func parse() (*Flags, error) {
	beforeDir := flag.String("before", "before", "Директория с исходными JSON файлами")
	afterDir := flag.String("after", "after", "Директория с измененными JSON файлами")
	showAll := flag.Bool("all", false, "Показать все файлы (даже без изменений)")
	workers := flag.Int("workers", runtime.NumCPU()-1, "Количество параллельных воркеров")
	flag.Parse()

	if *beforeDir == "" {
		fmt.Println(usage)
		return &Flags{}, fmt.Errorf("empty dir before %s", *beforeDir)
	}
	if *afterDir == "" {
		fmt.Println(usage)
		return &Flags{}, fmt.Errorf("empty dir after %s", *afterDir)
	}

	if *workers < 1 || *workers > runtime.NumCPU() {
		fmt.Println(usage)
		return &Flags{}, fmt.Errorf("workers=%d must be in [1, %d] ", *workers, runtime.NumCPU())
	}

	return &Flags{
		BeforeDir: *beforeDir,
		AfterDir:  *afterDir,
		ShowAll:   *showAll,
		Workers:   *workers,
	}, nil
}
