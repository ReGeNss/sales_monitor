package main

type Config struct {
	Timezone   string           `yaml:"timezone"`
	Shops      []ShopConfig     `yaml:"shops"`
	Categories []CategoryConfig `yaml:"categories"`
}

type ShopConfig struct {
	ID          string `yaml:"id"`
	DefaultCron string `yaml:"default_cron"`
}

type CategoryConfig struct {
	Name            string      `yaml:"name"`
	WordsToIgnore   []string    `yaml:"words_to_ignore"`
	Differentiation [][]string  `yaml:"differentiation"`
	Jobs            []JobConfig `yaml:"jobs"`
}

type JobConfig struct {
	ID           string   `yaml:"id"`
	ShopID       string   `yaml:"shop_id"`
	URLs         []string `yaml:"urls"`
	CronOverride string   `yaml:"cron_override"`
}

type ResolvedJob struct {
	ID              string
	ShopID          string
	Cron            string
	CategoryName    string
	WordsToIgnore   []string
	Differentiation [][]string
	URLs            []string
}
