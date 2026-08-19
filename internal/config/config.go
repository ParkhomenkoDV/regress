package config

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

type Config struct {
	Before  string `doc:"Директория ДО"`
	After   string `doc:"Директория ПОСЛЕ"`
	ShowAll bool   `doc:"Показывать все поля"`
	Workers int    `doc:"Количество параллельных работников"`
}

func New() (*Config, error) {
	flags, err := parse()
	if err != nil {
		return &Config{}, err
	}

	if flags.Before == "" {
		return &Config{}, fmt.Errorf("empty dir before %s", flags.Before)
	}
	if _, err := os.Stat(flags.Before); os.IsNotExist(err) {
		return &Config{}, fmt.Errorf("no found dir:%s", flags.Before)
	}
	if flags.After == "" {
		return &Config{}, fmt.Errorf("empty dir after %s", flags.After)
	}
	if _, err := os.Stat(flags.Before); os.IsNotExist(err) {
		return &Config{}, fmt.Errorf("no found dir:%s", flags.Before)
	}
	if flags.Workers < 1 {
		return &Config{}, fmt.Errorf("workers=%d must be >= 1", flags.Workers)
	}

	return flags, nil
}

func parse() (*Config, error) {
	before := flag.String(
		"before",
		"before",
		"Директория с исходными JSON файлами",
	)
	after := flag.String(
		"after",
		"after",
		"Директория с измененными JSON файлами",
	)
	showAll := flag.Bool(
		"all",
		false,
		"Показать все файлы (даже без изменений)",
	)
	workers := flag.Int(
		"workers",
		runtime.NumCPU(),
		"Количество параллельных воркеров",
	)

	flag.Parse()

	return &Config{
		Before:  *before,
		After:   *after,
		ShowAll: *showAll,
		Workers: *workers,
	}, nil
}
