package config

type Config struct {
	BeforeDir string
	AfterDir  string
	ShowAll   bool
	Workers   int
}

func New() (*Config, error) {
	flags, err := parse()
	if err != nil {
		return &Config{}, err
	}

	return &Config{
		BeforeDir: flags.BeforeDir,
		AfterDir:  flags.AfterDir,
		ShowAll:   flags.ShowAll,
		Workers:   flags.Workers,
	}, nil
}
