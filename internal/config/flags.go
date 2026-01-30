package config

import (
	"flag"
	"fmt"
	"runtime"
)

const usage = "Использование: go run regress.go [-before=S] [-after=S] [-all] [-workers=N]"

type Flags struct {
	BeforeDir string
	AfterDir  string
	ShowAll   bool
	Workers   int
}

func parse() (*Flags, error) {
	numCPU := runtime.NumCPU()

	beforeDir := flag.String("before", "before", "Директория с исходными JSON файлами")
	afterDir := flag.String("after", "after", "Директория с измененными JSON файлами")
	showAll := flag.Bool("all", false, "Показать все файлы (даже без изменений)")
	workers := flag.Int("workers", numCPU, "Количество параллельных воркеров")
	flag.Parse()

	if *beforeDir == "" {
		fmt.Println(usage)
		return &Flags{}, fmt.Errorf("empty dir before %s", *beforeDir)
	}
	if *afterDir == "" {
		fmt.Println(usage)
		return &Flags{}, fmt.Errorf("empty dir after %s", *afterDir)
	}

	if !(1 <= *workers && *workers <= numCPU) {
		fmt.Println(usage)
		return &Flags{}, fmt.Errorf("workers=%d must be in [1, %d] ", *workers, numCPU)
	}

	return &Flags{
		BeforeDir: *beforeDir,
		AfterDir:  *afterDir,
		ShowAll:   *showAll,
		Workers:   *workers,
	}, nil
}
